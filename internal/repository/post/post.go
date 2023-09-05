package post

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"forum/internal/types"
)

type PostDB struct {
	db *sql.DB
}

func NewPostDB(db *sql.DB) *PostDB {
	return &PostDB{db: db}
}

func (p *PostDB) GetAllPosts() ([]*types.Post, error) {
	// query := "SELECT * FROM snippets ORDER BY id DESC"

	// rows, err := p.db.Query(query)
	// if err != nil {
	// 	return nil, err
	// }
	// defer rows.Close()
	// var posts []*types.Post
	// for rows.Next() {
	// 	post := types.Post{}
	// 	err := rows.Scan(&post.Id, &post.AuthorId, &post.AuthorName, &post.Title, &post.Content, &post.Created)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	post.Time = post.Created.Format("15:04 January 02, 2006")
	// 	// get categories
	// 	rows1, err := p.db.Query("SELECT category FROM categories WHERE snippet_id= $1", post.Id)
	// 	if err != nil {
	// 		fmt.Println("GetAllPostsERR: ", err)
	// 	}
	// 	for rows1.Next() {
	// 		var category string
	// 		err := rows1.Scan(&category)
	// 		if err != nil {
	// 			return nil, err
	// 		}
	// 		post.Categories = append(post.Categories, category)
	// 	}
	// 	posts = append(posts, &post)
	// }
	// for _, post := range posts {
	// 	post.Likes = p.CountLikes(post.Id, "like")
	// 	post.Dislikes = p.CountLikes(post.Id, "dislike")
	// 	post.Comments = p.GetAllComments(post.Id)
	// }
	
	query := "SELECT id FROM snippets ORDER BY id DESC"
	rows, err := p.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []*types.Post
	for rows.Next() {
		post := &types.Post{}
		err = rows.Scan(
			&post.Id,
		)
		if err != nil {
			return nil, err
		}
		post, err = p.GetPostByID(post.Id)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}

func (p *PostDB) CreatePostDB(post *types.Post) (int, error) {
	res, err := p.db.Exec("INSERT INTO snippets (user_id, user_name, title, content, created) VALUES ($1, $2, $3, $4, DATETIME('now', '+6 hours'))",
		post.AuthorId,
		post.AuthorName,
		post.Title,
		post.Content,
	)
	if err != nil {
		fmt.Println("repository post err:", err)
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	if len(post.Categories) == 0 {
		_, err := p.db.Exec("INSERT INTO categories (snippet_id, category) VALUES ($1, $2)", id, "Other")
		if err != nil {
			fmt.Println("categories:", err)
			return 0, err
		}
	} else {
		for _, category := range post.Categories {
			_, err := p.db.Exec("INSERT INTO categories (snippet_id, category) VALUES ($1, $2)", id, category)
			if err != nil {
				fmt.Println("categories:", err)
				return 0, err
			}
		}
	}
	// Create image entry
	if post.ImageData != nil {
		_, err = p.db.Exec("INSERT INTO images (snippet_id, image, format) VALUES ($1, $2, $3)", id, post.ImageData, post.ImageFormat)
		if err != nil {
			return 0, err
		}
	}

	return int(id), nil
}

func (p *PostDB) GetPostByID(id int) (*types.Post, error) {
	post := &types.Post{}
	err := p.db.QueryRow("SELECT * FROM snippets WHERE id= $1", id).Scan(
		&post.Id,
		&post.AuthorId,
		&post.AuthorName,
		&post.Title,
		&post.Content,
		&post.Created)
	if err != nil {
		return nil, err
	}
	post.Time = post.Created.Format("15:04 January 02, 2006")
	rows, err := p.db.Query("SELECT category FROM categories WHERE snippet_id= $1", id)
	for rows.Next() {

		var category string
		err := rows.Scan(&category)
		if err != nil {
			return nil, err
		}
		post.Categories = append(post.Categories, category)
	}

	post.Likes = p.CountLikes(post.Id, "like")
	post.Dislikes = p.CountLikes(post.Id, "dislike")
	post.Comments = p.GetAllComments(post.Id)

	err = p.db.QueryRow("SELECT image, format FROM images WHERE snippet_id= $1", id).Scan(
		&post.ImageData,
		&post.ImageFormat,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			post.ImageData = nil
			err = nil
		} else {
			return nil, err
		}
	}
	post.ImageEncode = base64.StdEncoding.EncodeToString(post.ImageData)

	return post, err
}

func (p *PostDB) GetPostsByUserID(id int) ([]*types.Post, error) {
	rows, err := p.db.Query("SELECT id FROM snippets WHERE user_id = $1 ORDER BY id DESC", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []*types.Post
	for rows.Next() {
		post := &types.Post{}
		err := rows.Scan(
			&post.Id,
		)
		if err != nil {
			return nil, err
		}
		post, err = p.GetPostByID(post.Id)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	// for rows.Next() {
	// 	post := types.Post{}
	// 	err := rows.Scan(&post.Id, &post.AuthorId, &post.AuthorName, &post.Title, &post.Content, &post.Created)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	post.Time = post.Created.Format("15:04 January 02, 2006")

	// 	// get categories
	// 	rows1, err := p.db.Query("SELECT category FROM categories WHERE snippet_id= $1", post.Id)
	// 	if err != nil {
	// 		fmt.Println("GetAllPostsERR: ", err)
	// 	}
	// 	for rows1.Next() {

	// 		var category string
	// 		err := rows1.Scan(&category)
	// 		if err != nil {
	// 			return nil, err
	// 		}
	// 		post.Categories = append(post.Categories, category)
	// 	}

	// 	posts = append(posts, &post)

	// }
	// for _, post := range posts {
	// 	post.Likes = p.CountLikes(post.Id, "like")
	// 	post.Dislikes = p.CountLikes(post.Id, "dislike")
	// 	post.Comments = p.GetAllComments(post.Id)
	// }
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}
