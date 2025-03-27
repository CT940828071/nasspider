package provider

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"time"
	"strings"

	"github.com/chromedp/chromedp"
)

type DoMP4Provider struct{}

// ParseURLs 解析xpath获取URLs,当前集
func (d DoMP4Provider) ParseURLs(URL string, currentEp int) ([]string, int, error) {
	// 创建一个上下文
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-setuid-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// 设置超时时间
	ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	// 用于存储结果的变量
	var inputValues []string
	var titleText string

	// 运行任务
	err := chromedp.Run(ctx,
		// 访问目标网页
		chromedp.Navigate(URL),
		// 等待标题元素可见
		chromedp.WaitVisible(`/html/body/div[1]/div[2]/div[1]/h1`, chromedp.BySearch),
		// 获取标题元素的文本内容
		chromedp.Text(`/html/body/div[1]/div[2]/div[1]/h1`, &titleText, chromedp.BySearch),
		// 等待页面加载完成
		chromedp.WaitVisible(`#download1 ul li input`, chromedp.ByQuery),
		// 获取所有 input 元素的 value 属性
		chromedp.Evaluate(
			`Array.from(document.querySelectorAll('#download1 ul li input')).map(input => input.value);`,
			&inputValues,
		),
	)
	if err != nil {
		return nil, 0, err
	}

	// 定义正则表达式，匹配更新至全X集或更新至X集或更新至全X话或更新至X话
	re := regexp.MustCompile(`更新至(?:全)?(\d+)(?:集|话)`)
	matches := re.FindStringSubmatch(titleText)

	if len(matches) == 0 {
		// 再尝试解析全X集或X集或全X话或X话
		re = regexp.MustCompile(`(?:全)?(\d+)(?:集|话)`)
		matches = re.FindStringSubmatch(titleText)
		if len(matches) == 0 {
			return nil, 0, fmt.Errorf("异常:解析集数失败")
		}
	}

	updatedEp, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil, 0, fmt.Errorf("异常:集数转换失败")
	}

	// 解析集数小于当前记录集数
	if updatedEp < currentEp {
		return nil, 0, fmt.Errorf("异常：解析集数(%d) < 当前集数(%d)", updatedEp, currentEp)
	}

	if len(inputValues) != updatedEp {
		return nil, 0, fmt.Errorf("异常：解析集数(%d) != 下载链接数(%d)", updatedEp, len(inputValues))
	}

	// 解析集数等于当前记录集数
	if updatedEp == currentEp {
		return []string{}, 0, nil
	}

	// 根据匹配到的"集"或"话"来决定如何获取下载链接
	if strings.Contains(matches[0], "集") {
		// 匹配到"集"，从数组前面获取
		return inputValues[:updatedEp-currentEp], updatedEp, nil
	} else if strings.Contains(matches[0], "话") {
		// 匹配到"话"，从数组末尾获取
		startIndex := len(inputValues) - (updatedEp - currentEp)
		if startIndex < 0 {
			startIndex = 0
		}
		return inputValues[startIndex:], updatedEp, nil
	}

	return nil, 0, fmt.Errorf("异常:无法确定获取下载链接的方式")
}
