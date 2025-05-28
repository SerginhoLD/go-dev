package repository

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockHttpClient struct {
}

func (h *MockHttpClient) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString(MockPaginateResponse())),
	}, nil
}

func TestPaginate(t *testing.T) {
	t.Setenv("OBJECT_MAX_PRICE", "8900100")
	t.Setenv("OBJECT_DOMAIN", "https://site.local")
	t.Setenv("OBJECT_LIST", "/price-base-to-%[1]s/?PAGEN_1=%[2]s")
	t.Setenv("OBJECT_USER_AGENT", "Test")

	client := &MockHttpClient{}
	repository := NewVeObjectRepositoryImpl(client)
	objects := repository.Paginate(context.Background(), 1)

	assert.Equal(t, len(objects), 2)
}

func MockPaginateResponse() string {
	return `
<!DOCTYPE html>
<html lang="ru-RU">
<body>
<div class="container">
	<div class="title-listing">
		<h1>
			<span class="count">
				<span class="num">4083</span> 
				<span>объекта</span>
			</span>
		</h1>
	</div>
</div>

<div class="obj-list-page">
<div class="container">
<div class="obj-list">

<div class="obj-card">
	<a href="/objects/sale/flats/910895/" target="_blank">
					<div class="prices">
				<span class="price">9 500 000 &#8381;</span>
									<span class="price-permeter">256 065 &#8381;/м<sup>2</sup></span>
							</div>
				

		<div class="title">
			<h2>
									<span class="part">1 комн</span>
									<span class="part">37.1 м²</span>
									<span class="part">13/17 эт</span>
											</h2>
		</div>
	</a>
	<div class="address">г Москва, Варшавское шоссе, д 142 к 2</div>

                    <div class="metroblock">
                            <div class="metro">
                    <span class="txt">Пражская</span>
                    <span class="how foot" >6&nbsp;мин.</span>                </div>
                    </div>         
</div>

<div class="obj-card">
	<a href="/objects/sale/flats/923304/" target="_blank">
					<div class="soldout"><span>Объект продан</span></div>
		<div class="title">
			<h2>
									<span class="part">1 комн</span>
									<span class="part">37.8 м²</span>
									<span class="part">4/14 эт</span>
													<span class="obj-tt-small label-check" data-tooltip-content="#tt-obj-check"></span>
							</h2>
		</div>
	</a>
	<div class="address">г Москва, ул Омская, д 8</div>
                    <div class="metroblock">
                            <div class="metro">
                    <span class="txt">Аэропорт Внуково</span>
                    <span class="how transport" >7&nbsp;мин.</span>                </div>
                    </div>
            
</div>

</div>
</div>
</div>
</body>
</html>
	`
}
