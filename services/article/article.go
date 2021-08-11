package article

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	. "gen/models"
	"gen/registry"
	"gen/services/cache"
	"time"
)

type ArticleService struct {
	SQLStore *SQLService `inject:""`

	Cache *cache.CacheService `inject:""`
}

func init() {
	registry.RegisterService(&ArticleService{})
}

var (
	ErrorPermissionDenied = errors.New("没有操作权限")
)

func (r ArticleService) Init() error {
	return nil
}

func (r ArticleService) Create(param *CreateArticleCommand) (*Article, error) {
	article, err := CreateArticle(param)
	if err != nil {
		return nil, err
	}
	return article, nil
}

func (r ArticleService) Update(param *UpdateArticleCommand) error {
	article, err := GetArticleById(param.Id)
	if err != nil {
		return err
	}
	// 只有更新自己的文章
	if article.UserId != param.UserId {
		return ErrorPermissionDenied
	}
	if err := UpdateArticle(param); err != nil {
		return err
	}
	return nil
}

func (r ArticleService) GetById(id int) (*Article, error) {
	// 读取redis缓存
	c, err := r.Cache.Redis.Get(context.TODO(), fmt.Sprintf("article_id_%d", id)).Result()
	if err == nil {
		var a Article
		err := json.Unmarshal([]byte(c), &a)
		if err == nil {
			err = AddViewNum(id)
			if err != nil {
				return nil, err
			}
			return &a, nil
		}
	}
	article, err := GetArticleById(id)
	if err != nil {
		return nil, err
	}
	// 缓存数据
	data, err := json.Marshal(article)
	if err == nil {
		r.Cache.Redis.Set(context.TODO(), fmt.Sprintf("article_id_%d", id), data, time.Second*30)
	}
	err = AddViewNum(id)
	if err != nil {
		return nil, err
	}
	return article, nil
}

func (r ArticleService) GetAll(page int) ([]*Article, error) {
	const pageSize = 15
	articles, err := GetArticles(page, pageSize)
	if err != nil {
		return nil, err
	}
	return articles, nil
}

func (r ArticleService) Delete(id, userId int) error {
	article, err := GetArticleById(id)
	if err != nil {
		return err
	}
	if article.UserId != userId {
		return ErrorPermissionDenied
	}
	if err = DeleteArticle(id); err != nil {
		return err
	}
	return nil
}

func (r ArticleService) AddComment(param *CreateArticleCommentCommand) error {
	_, err := GetArticleById(param.Id)
	if err != nil {
		return err
	}
	if err := CreateArticleComment(param); err != nil {
		return err
	}
	return nil
}
