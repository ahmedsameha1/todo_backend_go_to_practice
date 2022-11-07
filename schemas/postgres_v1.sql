create table todo (
    id uuid primary key,
    title varchar(500) not null,
    description varchar(10000) not null,
    done bool not null,
    created_at timestamptz not null,
    user_id varchar(40) not null
);
