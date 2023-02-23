CREATE TABLE IF NOT EXISTS books (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    title text NOT NULL,
    authors text NOT NULL,
    rating numeric(3,2) NOT NULL,
    ISBN text not null ,
    ISBN13 text not null,
    language varchar(20) not null,
    pages integer NOT NULL,
    genres text[] NOT NULL,
    version integer NOT NULL DEFAULT 1
);