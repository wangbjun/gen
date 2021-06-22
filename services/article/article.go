package article

import (
	"errors"
	. "gen/models"
	"gen/registry"
)

type ArticleService struct {
	SQLStore *SQLStore `inject:""`
}

func init() {
	registry.RegisterService(&ArticleService{})
}

var (
	PermissionDenied = errors.New("没有操作权限")
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
		return PermissionDenied
	}
	if err := UpdateArticle(param); err != nil {
		return err
	}
	return nil
}

func (r ArticleService) GetById(id int) (*Article, error) {
	article, err := GetArticleById(id)
	if err != nil {
		return nil, err
	}
	err = IncArticleViewNum(id)
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
		return PermissionDenied
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
