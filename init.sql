create extension citext;

create table test_results(
    id serial,
    u_id int,
    ex_id int,
    mark int,
    error citext
    compiled bool,
    compile_output citext,
    file_name citext,
    file_name_old citext,
    run_output jsonb
)