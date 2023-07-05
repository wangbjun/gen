package services

import (
	"context"
	"gen/log"
	. "gen/models"
)

type ArticleService struct{}

func NewArticleService() *ArticleService {
	return &ArticleService{}
}

func (r ArticleService) Create(ctx context.Context, param *CreateArticleCommand) (*Article, error) {
	article, err := CreateArticle(ctx, param)
	if err != nil {
		return nil, err
	}
	return article, nil
}

func (r ArticleService) Update(ctx context.Context, param *UpdateArticleCommand) error {
	_, err := GetArticleById(ctx, param.Id)
	if err != nil {
		return err
	}
	if err := UpdateArticle(ctx, param); err != nil {
		return err
	}
	return nil
}

func (r ArticleService) GetById(ctx context.Context, id int) (*Article, error) {
	article, err := GetArticleById(ctx, id)
	if err != nil {
		return nil, err
	}
	go func() {
		err = AddViewNum(ctx, id)
		if err != nil {
			log.WithCtx(ctx).Error("AddViewNum failed: %d", id)
		}
	}()
	return article, nil
}

func (r ArticleService) GetAll(ctx context.Context, page, pageSize int) ([]*Article, int64, error) {
	articles, totalCount, err := GetArticles(ctx, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return articles, totalCount, nil
}

func (r ArticleService) Delete(ctx context.Context, id int) error {
	_, err := GetArticleById(ctx, id)
	if err != nil {
		return err
	}
	if err = DeleteArticle(ctx, id); err != nil {
		return err
	}
	return nil
}

func (r ArticleService) AddComment(ctx context.Context, param *CreateArticleCommentCommand) error {
	_, err := GetArticleById(ctx, param.Id)
	if err != nil {
		return err
	}
	if err := CreateArticleComment(ctx, param); err != nil {
		return err
	}
	return nil
}
