package repository

import (
	"context"
	"exampleapp/internal/domain/entity"
	"exampleapp/internal/domain/repository"
	"exampleapp/internal/infrastructure/manticore"
	"slices"
	"time"
)

type ObjectRepositoryImpl struct {
	client *manticore.Client
}

func NewObjectRepositoryImpl(client *manticore.Client) repository.ObjectRepository {
	return &ObjectRepositoryImpl{client}
}

func (r *ObjectRepositoryImpl) Paginate(ctx context.Context, query repository.ObjectsQuery) (objects []*entity.Object, total uint64) {
	search := r.client.NewSearch("objects", query.Page, query.Limit)

	if query.Metro != "" {
		search = search.Equals("metro", query.Metro)
	}

	if query.Checked == -1 {
		search = search.Equals("checked", 0)
	} else if query.Checked == 1 {
		search = search.Equals("checked", 1)
	}

	if len(query.Search) > 0 {
		search = search.Match("_all", query.Search)
	}

	if query.PriceFrom > 0 {
		search.Range("price", "gte", query.PriceFrom)
	}

	if query.PriceTo > 0 {
		search.Range("price", "lte", query.PriceTo)
	}

	if query.SizeFrom > 0 {
		search.Range("size", "gte", query.SizeFrom)
	}

	if query.SizeTo > 0 {
		search.Range("size", "lte", query.SizeTo)
	}

	if query.Rooms > 0 {
		search = search.Equals("rooms", query.Rooms)
	}

	if query.Loc == 1 {
		search = search.In("metro", r.locCenter()...)
	} else if query.Loc == 2 {
		search = search.NotIn("metro", r.locCenter()...)
	}

	response, err := search.Execute(ctx)

	if err != nil {
		panic(err.Error())
	}

	if *response.Hits.Total == 0 {
		//return []*entity.Object{}, 0
		return objects, 0
	}

	//var objects []*entity.Object

	for _, hit := range response.Hits.GetHits() {
		source := hit.GetSource()

		objects = append(objects, &entity.Object{
			Id:        uint64(hit.GetId()),
			Title:     source["title"].(string),
			Metro:     source["metro"].(string),
			Price:     uint64(source["price"].(float64)),
			Size:      source["size"].(float64),
			Rooms:     uint8(source["rooms"].(float64)),
			Checked:   source["checked"].(bool),
			UpdatedAt: time.Unix(int64(source["updated_at"].(float64)), 0),
		})
	}

	return objects, uint64(*response.Hits.Total)
}

func (r *ObjectRepositoryImpl) Save(ctx context.Context, objects []*entity.Object) {
	for _, obj := range objects {
		r.client.Replace(ctx, "objects", obj.Id, map[string]any{
			"title":      obj.Title,
			"metro":      obj.Metro,
			"price":      obj.Price,
			"size":       obj.Size,
			"rooms":      obj.Rooms,
			"checked":    obj.Checked,
			"updated_at": obj.UpdatedAt,
		})
	}
}

func (r *ObjectRepositoryImpl) Metro(ctx context.Context) []string {
	buckets, err := r.client.AggsTerm(ctx, "objects", "metro")

	if err != nil {
		panic(err.Error())
	}

	var metro []string

	for _, bucket := range buckets {
		metro = append(metro, bucket.Key.(string))
	}

	slices.Sort(metro)
	return metro
}

func (r *ObjectRepositoryImpl) Rooms(ctx context.Context) []uint8 {
	buckets, err := r.client.AggsTerm(ctx, "objects", "rooms")

	if err != nil {
		panic(err.Error())
	}

	var rooms []uint8

	for _, bucket := range buckets {
		rooms = append(rooms, uint8(bucket.Key.(float64)))
	}

	slices.Sort(rooms)
	return rooms
}

func (r *ObjectRepositoryImpl) locCenter() []any {
	return []any{
		"Тульская", "Серпуховская", "Добрынинская", "Полянка", "Боровицкая", "Чеховская", "Цветной бульвар", "Менделевская", "Савёловская", "Новослободская",
		"Марьина Роща", "Достоевская", "Трубная", "Сретенский бульвар", "Чкаловская", "Римская", "Крестьянская Застава", "Дубровка", "Кожуховская",
		"Рижская", "Проспект Мира", "Сухаревская", "Тургеневская", "Китай-город", "Третьяковская", "Октябрьская", "Шаболовская", "Ленинский проспект", "Площадь Гагарина",
		"Сокольники", "Митьково", "Красносельская", "Площадь трёх вокзалов", "Ленинградский вокзал", "Казанский вокзал", "Комсомольская", "Красные Ворота", "Чистые пруды", "Лубянка", "Охотный Ряд", "Библиотека имени Ленина", "Кропоткинская", "Парк культуры", "Фрунзенская", "Спортивная", "Лужники", "Воробьёвы горы", "Университет",
		"Измайловская", "Измайлово", "Партизанская", "Соколиная Гора", "Семёновская", "Электрозаводская", "Сортировочная", "Лефортово", "Бауманская", "Курская", "Площадь Революции", "Арбатская", "Смоленская", "Киевская", "Студенческая", "Деловой центр",
		"Шоссе Энтузиастов", "Авиамоторная", "Андроновка", "Нижегородская", "Москва-Товарная", "Площадь Ильича", "Серп и Молот", "Калитники", "Марксистская", "Таганская", "Стахановская",
		"Текстильщики", "Волгоградский проспект", "Пролетарская", "Кузнецкий Мост", "Пушкинская", "Баррикадная", "Краснопресненская", "Улица 1905 года", "Беговая",
		"Автозаводская", "ЗИЛ", "Павелецкая", "Новокузнецкая", "Театральная", "Тверская", "Маяковская", "Белорусская", "Петровский парк", "Динамо", "ЦСКА", "Зорге",
	}
}
