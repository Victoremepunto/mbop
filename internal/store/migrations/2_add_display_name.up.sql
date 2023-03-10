alter table registrations
    add display_name varchar default NULL;

-- backfill all null display names
do
$$
    declare
        row RECORD;
    begin
        for row in
            select * from registrations where display_name is null
            loop
                update registrations set display_name = uid where id = row.id;
            end loop;
    end
$$;

create unique index display_name_unique
    on registrations (display_name);
