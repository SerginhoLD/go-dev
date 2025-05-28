package repository

import (
	"context"
	"exampleapp/internal/domain/entity"
	"exampleapp/internal/domain/repository"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type VeObjectRepositoryImpl struct {
	client HttpClient
}

func NewVeObjectRepositoryImpl(client HttpClient) repository.VeObjectRepository {
	return &VeObjectRepositoryImpl{client}
}

func (r *VeObjectRepositoryImpl) Total(ctx context.Context) (total uint64, size uint8) {
	url := fmt.Sprintf(os.Getenv("OBJECT_LIST"), os.Getenv("OBJECT_MAX_PRICE"), "1")
	url = fmt.Sprintf("%s%s", os.Getenv("OBJECT_DOMAIN"), url)
	slog.DebugContext(ctx, fmt.Sprintf(`http: "%s"`, url))

	req := NewHttpRequest("GET", url)
	resp, err := r.client.Do(req)

	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf(`http: %s`, err.Error()))
		return 0, 0
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		slog.ErrorContext(ctx, fmt.Sprintf(`http: "%s" %d`, url, resp.StatusCode))
		return 0, 0
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf(`http: %s`, err.Error()))
		return 0, 0
	}

	total, err = strconv.ParseUint(doc.Find(".title-listing .num").Text(), 10, 64)
	size = uint8(doc.Find(".obj-card").Length())

	if err != nil || size == 0 {
		return 0, 0
	}

	slog.DebugContext(ctx, fmt.Sprintf(`StartParseObjects: total="%d" size="%d"`, total, size))
	return total, size
}

func (r *VeObjectRepositoryImpl) Paginate(ctx context.Context, page uint64) []*entity.Object {
	url := fmt.Sprintf(os.Getenv("OBJECT_LIST"), os.Getenv("OBJECT_MAX_PRICE"), strconv.FormatUint(page, 10))
	url = fmt.Sprintf("%s%s", os.Getenv("OBJECT_DOMAIN"), url)
	slog.DebugContext(ctx, fmt.Sprintf(`http: "%s"`, url))

	req := NewHttpRequest("GET", url)
	resp, err := r.client.Do(req)

	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf(`http: %s`, err.Error()))
		return []*entity.Object{}
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		slog.ErrorContext(ctx, fmt.Sprintf(`http: "%s" %d`, url, resp.StatusCode))
		return []*entity.Object{}
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		//panic(err)
		slog.ErrorContext(ctx, fmt.Sprintf(`http: %s`, err.Error()))
		return []*entity.Object{}
	}

	total, err := strconv.ParseUint(doc.Find(".title-listing .num").Text(), 10, 64)

	if err != nil || total == 0 {
		return []*entity.Object{}
	}

	var objects []*entity.Object

	doc.Find(".obj-card").Each(func(i int, seCard *goquery.Selection) {
		title := seCard.Find(".address").Text()
		size := float64(0)
		rooms := uint8(0)

		seCard.Find(".title .part").Each(func(i int, sePart *goquery.Selection) {
			title = fmt.Sprintf("%s, %s", title, sePart.Text())

			if i == 0 {
				reRooms := regexp.MustCompile(`\d+`)
				rooms64, _ := strconv.ParseUint(reRooms.FindString(sePart.Text()), 10, 64)
				rooms = uint8(rooms64)
			}

			if i == 1 {
				reSize := regexp.MustCompile(`\d+\.?\d+?`)
				size, _ = strconv.ParseFloat(reSize.FindString(sePart.Text()), 64)
			}
		})

		reSlug := regexp.MustCompile(`\d+`)
		slug := reSlug.FindString(seCard.Find("a").AttrOr("href", ""))

		if len(slug) != 0 {
			id, _ := strconv.ParseUint(slug, 10, 64)
			price := uint64(0)

			if seCard.Find(".price").Length() > 0 {
				rePrice := regexp.MustCompile(`^(\d+)\s?(\d+)?\s?(\d+)?.*`)
				reMatch := rePrice.FindAllStringSubmatch(seCard.Find(".price").Text(), -1)
				price, _ = strconv.ParseUint(fmt.Sprintf("%s%s%s", reMatch[0][1], reMatch[0][2], reMatch[0][3]), 10, 64)
			}

			objects = append(objects, &entity.Object{
				Id:        id,
				Title:     title,
				Metro:     seCard.Find(".txt").Text(),
				Price:     price,
				Size:      size,
				Rooms:     rooms,
				Checked:   seCard.Find(".label-check").Length() > 0,
				UpdatedAt: time.Now(),
			})
		}
	})

	return objects
}
