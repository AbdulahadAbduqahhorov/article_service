BEGIN;

INSERT INTO author (id,firstname, lastname) VALUES ('de9fc112-7511-4abd-9222-7019b65d1108', 'John', 'Smith') ON CONFLICT DO NOTHING;
INSERT INTO author (id,firstname, lastname) VALUES ( '9ae20328-995c-401b-90a0-9acadf63c6ec', 'Abdulahad', 'Abduqahhorov') ON CONFLICT DO NOTHING;
INSERT INTO article (id,title, body,author_id) VALUES ('ead02d29-5bf9-4b9e-91c0-c6e1ab648937', 'title 1', 'body 1','9ae20328-995c-401b-90a0-9acadf63c6ec') ON CONFLICT DO NOTHING;
INSERT INTO article (id,title, body,author_id) VALUES ('0bd2034c-5283-4e59-b904-02ff9fa8ed48', 'title 2','body 2','9ae20328-995c-401b-90a0-9acadf63c6ec') ON CONFLICT DO NOTHING;
INSERT INTO article (id,title, body,author_id) VALUES ('cb7cfd12-6501-4e8c-8297-cb442512b9ba', 'title 3', 'body 3', 'de9fc112-7511-4abd-9222-7019b65d1108') ON CONFLICT DO NOTHING;
INSERT INTO article (id,title, body,author_id) VALUES ('0579a273-a581-4806-8e0b-4b9188c852ec', 'title 4', 'body 4','de9fc112-7511-4abd-9222-7019b65d1108') ON CONFLICT DO NOTHING;

COMMIT;