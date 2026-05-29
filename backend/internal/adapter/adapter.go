package adapter

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
	"time"

	"content-publisher/backend/internal/domain"
)

type PlatformAdapter interface {
	Descriptor() domain.PlatformDescriptor
	BuildPreview(content domain.ContentInput) domain.PlatformPreview
	Publish(payload domain.PublishPayload) domain.PublishResult
}

type Registry struct {
	adapters map[string]PlatformAdapter
}

func NewRegistry(adapters ...PlatformAdapter) *Registry {
	items := make(map[string]PlatformAdapter, len(adapters))
	for _, item := range adapters {
		items[item.Descriptor().ID] = item
	}
	return &Registry{adapters: items}
}

func (r *Registry) List() []domain.PlatformDescriptor {
	descriptors := make([]domain.PlatformDescriptor, 0, len(r.adapters))
	for _, item := range r.adapters {
		descriptors = append(descriptors, item.Descriptor())
	}
	return descriptors
}

func (r *Registry) Get(id string) (PlatformAdapter, bool) {
	adapter, ok := r.adapters[id]
	return adapter, ok
}

type baseAdapter struct {
	id          string
	name        string
	description string
	styleHints  []string
	supports    []string
	titleLimit  int
	tagLimit    int
	bodyBuilder func(domain.ContentInput) string
	notes       func(domain.ContentInput) []string
	warnings    func(domain.ContentInput) []string
}

func (b baseAdapter) Descriptor() domain.PlatformDescriptor {
	return domain.PlatformDescriptor{
		ID:          b.id,
		Name:        b.name,
		Description: b.description,
		StyleHints:  b.styleHints,
		Supports:    b.supports,
	}
}

func (b baseAdapter) BuildPreview(content domain.ContentInput) domain.PlatformPreview {
	title := trimRunes(content.Title, b.titleLimit)
	tags := pickTags(content.Tags, b.tagLimit)

	return domain.PlatformPreview{
		PlatformID:      b.id,
		PlatformName:    b.name,
		Title:           title,
		FormattedBody:   b.bodyBuilder(content),
		RecommendedTags: tags,
		Notes:           b.notes(content),
		Warnings:        b.warnings(content),
	}
}

func (b baseAdapter) Publish(payload domain.PublishPayload) domain.PublishResult {
	now := time.Now()
	status := "demo_ready"
	message := "Demo mode completed. Real platform publishing is not connected yet."

	if payload.Simulate {
		status = "simulated"
		message = "Simulation completed. No external platform API was called."
	}

	hash := sha1.Sum([]byte(fmt.Sprintf("%s|%s|%s", b.id, payload.Content.Title, now.Format(time.RFC3339Nano))))

	return domain.PublishResult{
		PlatformID:   b.id,
		PlatformName: b.name,
		Status:       status,
		Message:      message,
		PublishedAt:  now,
		ExternalRef:  strings.ToUpper(b.id) + "-" + hex.EncodeToString(hash[:])[:10],
		Simulation:   payload.Simulate,
	}
}

func NewWeChatOfficialAccountAdapter() PlatformAdapter {
	return baseAdapter{
		id:          "wechat",
		name:        "公众号",
		description: "长文图文内容，强调章节结构、摘要与引导关注。",
		styleHints:  []string{"章节分明", "适合长文", "结尾引导互动"},
		supports:    []string{"rich_text", "cover_image", "scheduled_publish"},
		titleLimit:  64,
		tagLimit:    5,
		bodyBuilder: func(content domain.ContentInput) string {
			var parts []string
			parts = append(parts, "# "+content.Title)
			if strings.TrimSpace(content.Summary) != "" {
				parts = append(parts, "> 摘要："+content.Summary)
			}
			parts = append(parts, normalizeParagraphs(content.Body, "\n\n"))
			parts = append(parts, "\n---\n欢迎在评论区交流，关注我获取更多内容。")
			return strings.Join(parts, "\n\n")
		},
		notes: func(content domain.ContentInput) []string {
			notes := []string{"建议上传 2.35:1 封面图，增强文章点击率。"}
			if strings.TrimSpace(content.Summary) == "" {
				notes = append(notes, "建议补充摘要，便于生成公众号导语。")
			}
			return notes
		},
		warnings: func(content domain.ContentInput) []string {
			warnings := []string{}
			if len([]rune(content.Body)) < 180 {
				warnings = append(warnings, "正文偏短，公众号读者通常期待更完整的阐述。")
			}
			return warnings
		},
	}
}

func NewZhihuAdapter() PlatformAdapter {
	return baseAdapter{
		id:          "zhihu",
		name:        "知乎",
		description: "知识型内容，强调观点、论证与小标题。",
		styleHints:  []string{"观点先行", "论证充分", "适合问答和专栏"},
		supports:    []string{"markdown", "topics", "scheduled_publish"},
		titleLimit:  80,
		tagLimit:    6,
		bodyBuilder: func(content domain.ContentInput) string {
			body := normalizeParagraphs(content.Body, "\n\n")
			body = emphasizeQuestions(body)
			return fmt.Sprintf("## 核心观点\n\n%s\n\n## 展开说明\n\n%s\n\n## 总结\n\n%s", fallbackSummary(content), body, closingSentence(content))
		},
		notes: func(content domain.ContentInput) []string {
			return []string{"建议在开头 3 句内明确结论，提升知乎完读率。"}
		},
		warnings: func(content domain.ContentInput) []string {
			if len(content.Tags) < 2 {
				return []string{"标签较少，建议补充 2-3 个相关话题。"}
			}
			return []string{}
		},
	}
}

func NewBilibiliAdapter() PlatformAdapter {
	return baseAdapter{
		id:          "bilibili",
		name:        "B站动态",
		description: "适合视频配文和动态介绍，强调节奏和行动号召。",
		styleHints:  []string{"节奏紧凑", "适合视频导流", "可以搭配分点说明"},
		supports:    []string{"dynamic_post", "video_link", "topics"},
		titleLimit:  40,
		tagLimit:    4,
		bodyBuilder: func(content domain.ContentInput) string {
			body := normalizeParagraphs(content.Body, "\n")
			lines := strings.Split(body, "\n")
			condensed := make([]string, 0, len(lines))
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line == "" {
					continue
				}
				condensed = append(condensed, "• "+line)
			}
			condensed = append(condensed, "想看完整内容或演示，欢迎三连支持。")
			return strings.Join(condensed, "\n")
		},
		notes: func(content domain.ContentInput) []string {
			return []string{"建议补充视频链接或封面截图，提升动态点击率。"}
		},
		warnings: func(content domain.ContentInput) []string {
			if len([]rune(content.Body)) > 500 {
				return []string{"文案偏长，建议压缩为高信息密度短段落。"}
			}
			return []string{}
		},
	}
}

func NewXiaohongshuAdapter() PlatformAdapter {
	return baseAdapter{
		id:          "xiaohongshu",
		name:        "小红书",
		description: "种草与经验分享导向，强调生活化语气和标签曝光。",
		styleHints:  []string{"口语化", "强调体验", "结尾带话题标签"},
		supports:    []string{"note_post", "topics", "cover_image"},
		titleLimit:  20,
		tagLimit:    8,
		bodyBuilder: func(content domain.ContentInput) string {
			body := normalizeParagraphs(content.Body, "\n\n")
			tags := pickTags(content.Tags, 8)
			tagParts := make([]string, 0, len(tags))
			for _, tag := range tags {
				tagParts = append(tagParts, "#"+sanitizeTag(tag))
			}
			return fmt.Sprintf("%s\n\n%s\n\n%s", casualLead(content), body, strings.Join(tagParts, " "))
		},
		notes: func(content domain.ContentInput) []string {
			return []string{"建议搭配 3-9 张图片或短视频封面，增强种草效果。"}
		},
		warnings: func(content domain.ContentInput) []string {
			warnings := []string{}
			if strings.TrimSpace(content.CoverImage) == "" {
				warnings = append(warnings, "未设置封面图，小红书内容通常依赖视觉首图。")
			}
			return warnings
		},
	}
}

func trimRunes(input string, limit int) string {
	if limit <= 0 {
		return input
	}
	runes := []rune(strings.TrimSpace(input))
	if len(runes) <= limit {
		return string(runes)
	}
	return string(runes[:limit-1]) + "…"
}

func pickTags(tags []string, limit int) []string {
	if len(tags) == 0 {
		return []string{"内容创作", "效率工具"}
	}
	if len(tags) > limit {
		return append([]string{}, tags[:limit]...)
	}
	return append([]string{}, tags...)
}

func normalizeParagraphs(body string, separator string) string {
	parts := strings.Split(body, "\n")
	cleaned := make([]string, 0, len(parts))
	for _, part := range parts {
		line := strings.TrimSpace(part)
		if line != "" {
			cleaned = append(cleaned, line)
		}
	}
	return strings.Join(cleaned, separator)
}

func fallbackSummary(content domain.ContentInput) string {
	if strings.TrimSpace(content.Summary) != "" {
		return content.Summary
	}
	body := normalizeParagraphs(content.Body, " ")
	return trimRunes(body, 80)
}

func closingSentence(content domain.ContentInput) string {
	if content.Tone == "professional" {
		return "如果你也在优化内容分发流程，欢迎交流你的方法论。"
	}
	return "如果这篇内容对你有帮助，欢迎留言说说你的看法。"
}

func casualLead(content domain.ContentInput) string {
	if strings.TrimSpace(content.Summary) != "" {
		return "今天想分享一个很实用的经验：" + content.Summary
	}
	return "今天整理了一些亲测有效的内容发布思路，直接上干货。"
}

func sanitizeTag(tag string) string {
	tag = strings.TrimSpace(tag)
	tag = strings.TrimPrefix(tag, "#")
	return strings.ReplaceAll(tag, " ", "")
}

func emphasizeQuestions(body string) string {
	re := regexp.MustCompile(`([?？])`)
	if re.MatchString(body) {
		return body
	}
	return body
}
