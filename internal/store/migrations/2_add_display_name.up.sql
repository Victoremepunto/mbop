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

-- unique across multiple tenants per display name
alter table registrations
    add constraint display_name_unique
        unique (display_name);
