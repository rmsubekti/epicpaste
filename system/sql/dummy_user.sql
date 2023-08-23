--- insert user account example ---
insert into "user".account(id, user_name, email, password) 
    values ('2539b6ba-7ff6-4864-b7bc-6fed752ee925', 'epicpaster', 'epicpaster@epicpaste.com', '$2a$10$a99LcpBKvlekF17zvS63S.tHmLrA9wdVJnBAgOrpUmUs3N07pW2D2') on conflict do nothing;
    -- pasword is 5uperSecret
--- add user detail for created account
insert into "user".user(id, name) values ('epicpaster', 'Epic Paster') on conflict do nothing;


--- ===================================================================================================================================================== ---
--- Write more sql here ---