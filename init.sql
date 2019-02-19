create extension citext;

create table test_results(
    id serial,
    u_id int,
    ex_id int,
    mark int,
    error citext
)