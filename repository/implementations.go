package repository

import (
	"context"
	"fmt"
	"time"
)

func (r *Repository) GetTestById(ctx context.Context, input GetTestByIdInput) (output GetTestByIdOutput, err error) {
	err = r.Db.QueryRowContext(ctx, "SELECT name FROM test WHERE id = $1", input.Id).Scan(&output.Name)
	if err != nil {
		panic(err)
	}
	return
}

func (r *Repository) InsertNewUser(ctx context.Context, input UsersEntity) (output UsersEntity, err error) {
	err = r.Db.QueryRowContext(ctx, `INSERT INTO users (phone_number, full_name, passwd, success_login_count, created_at)
	VALUES($1, $2, $3, $4, $5) RETURNING id`, input.PhoneNumber, input.FullName, input.Password, 0, time.Now().UTC()).Scan(&output.Id)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	return
}

func (r *Repository) Login(ctx context.Context, input UsersEntity) (output UsersEntity, err error) {
	err = r.Db.QueryRowContext(ctx, `SELECT id, passwd
	FROM users WHERE phone_number = $1`, input.PhoneNumber).Scan(&output.Id, &output.Password)
	if err != nil {
		return
	}
	return
}

func (r *Repository) GetMyMyProfile(ctx context.Context, input UsersEntity) (output UsersEntity, err error) {
	err = r.Db.QueryRowContext(ctx, `SELECT phone_number, full_name
	FROM users WHERE id = $1`, input.Id).Scan(&output.PhoneNumber, &output.FullName)
	if err != nil {
		panic(err)
	}
	return
}

func (r *Repository) UpdateUser(ctx context.Context, input UsersEntity) (output UsersEntity, err error) {
	_, err = r.Db.ExecContext(ctx, `UPDATE users
	SET phone_number=$2, full_name=$3
	WHERE id=$1`, input.Id, input.PhoneNumber, input.FullName)
	if err != nil {
		panic(err)
	}
	return
}

func (r *Repository) UpdateLoginCount(ctx context.Context, input UsersEntity) (output UsersEntity, err error) {
	_, err = r.Db.ExecContext(ctx, `UPDATE users
	SET  success_login_count = 1 + success_login_count
	WHERE id=$1`, input.Id)
	if err != nil {
		panic(err)
	}
	return
}
