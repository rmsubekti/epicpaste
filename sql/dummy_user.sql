--- insert user account example ADMIN USER ---
insert into "user".account(id, user_name, email, password) 
    values ('2539b6ba-7ff6-4864-b7bc-6fed752ee925', 'admin', 'admin@epicpaste.com', '$2a$10$a99LcpBKvlekF17zvS63S.tHmLrA9wdVJnBAgOrpUmUs3N07pW2D2') on conflict do nothing;
    -- pasword is 5uperSecret
--- add user detail for created account
insert into "user".user(id, name) values ('2539b6ba-7ff6-4864-b7bc-6fed752ee925', 'User Admin') on conflict do nothing;


--- ===================================================================================================================================================== ---
--- Write more sql here ---