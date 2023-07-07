package service

import (
	"context"
	"fmt"
	"gen/log"
	. "gen/model"
)

type ArticleService struct{}

func NewArticleService() *ArticleService {
	return &ArticleService{}
}

func (r ArticleService) GetAll(ctx context.Context, page, pageSize int) ([]*Article, int64, error) {
	articles, totalCount, err := ArticleModel.GetAll(ctx, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return articles, totalCount, nil
}

func (r ArticleService) GetById(ctx context.Context, id int) (*Article, error) {
	article, err := ArticleModel.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	go func() {
		err = ArticleModel.AddViewNum(ctx, id)
		if err != nil {
			log.WithCtx(ctx).Error(fmt.Sprintf("AddViewNum failed: %d", id))
		}
	}()
	return article, nil
}

func (r ArticleService) Create(ctx context.Context, param *CreateArticleCommand) (*Article, error) {
	article, err := ArticleModel.Create(ctx, param)
	if err != nil {
		return nil, err
	}
	return article, nil
}

func (r ArticleService) Update(ctx context.Context, param *UpdateArticleCommand) error {
	_, err := ArticleModel.GetById(ctx, param.Id)
	if err != nil {
		return err
	}
	if err := ArticleModel.Update(ctx, param); err != nil {
		return err
	}
	return nil
}

func (r ArticleService) Delete(ctx context.Context, id int) error {
	_, err := ArticleModel.GetById(ctx, id)
	if err != nil {
		return err
	}
	if err = ArticleModel.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}

func (r ArticleService) AddComment(ctx context.Context, param *CreateArticleCommentCommand) error {
	_, err := ArticleModel.GetById(ctx, param.Id)
	if err != nil {
		return err
	}
	if err := ArticleModel.CreateComment(ctx, param); err != nil {
		return err
	}
	return nil
}
