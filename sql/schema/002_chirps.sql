-- +goose Up
CREATE TABLE chirps (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL ,
    updated_at TIMESTAMP NOT NULL,
    body TEXT NOT NULL,
    user_id UUID NOT NULL,
    constraint fk_user foreign key (user_id) references users(id) on delete cascade
);

-- +goose Down
DROP TABLE chirps;