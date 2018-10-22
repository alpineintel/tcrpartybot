----
-- Accounts
--
CREATE TABLE public.accounts (
    id                               SERIAL PRIMARY KEY NOT NULL,
    twitter_handle                   VARCHAR(15) NOT NULL,
    twitter_id                       BIGINT NOT NULL,
    eth_private_key                  VARCHAR(64) NOT NULL,
    passed_registration_challenge_at TIMESTAMP WITHOUT TIME ZONE,
    created_at                       TIMESTAMP WITHOUT TIME ZONE NOT NULL
        DEFAULT (now() AT TIME ZONE 'utc')
);

----
-- Registration questions
--
CREATE TABLE public.registration_questions (
    id            SERIAL PRIMARY KEY NOT NULL,
    question      CHARACTER VARYING NOT NULL,
    answer        CHARACTER VARYING NOT NULL
);

INSERT INTO public.registration_questions (question, answer) VALUES (
    'What is 2+2?',
    '4'
);

INSERT INTO public.registration_questions (question, answer) VALUES (
    'What is 4+1?',
    '5'
);

INSERT INTO public.registration_questions (question, answer) VALUES (
    'What is 10*5?',
    '50'
);

----
-- Registration challenges
--
CREATE TABLE public.registration_challenges (
    id                                    SERIAL PRIMARY KEY NOT NULL,
    account_id                            INTEGER NOT NULL,
    registration_question_id              INTEGER NOT NULL,
    sent_at                               TIMESTAMP WITHOUT TIME ZONE,
    completed_at                          TIMESTAMP WITHOUT TIME ZONE
);

----
-- OAuth Tokens
--
CREATE TABLE public.oauth_tokens (
    id                 SERIAL PRIMARY KEY NOT NULL,
    twitter_id         BIGINT NOT NULL,
    twitter_handle     VARCHAR(15) NOT NULL,
    oauth_token        VARCHAR(64) NOT NULL,
    oauth_token_secret VARCHAR(64) NOT NULL,
    created_at         TIMESTAMP WITHOUT TIME ZONE
        DEFAULT (now() AT TIME ZONE 'utc')
);

----
-- Key value store
--
CREATE TABLE public.keyval_store (
    key   CHARACTER VARYING NOT NULL PRIMARY KEY,
    value CHARACTER VARYING
);


----
-- Indices and such
--
ALTER TABLE ONLY public.registration_challenges
    ADD CONSTRAINT challenge_question FOREIGN KEY (registration_question_id)
    REFERENCES public.registration_questions(id);
ALTER TABLE ONLY public.registration_challenges
    ADD CONSTRAINT challenge_account FOREIGN KEY (account_id)
    REFERENCES public.accounts(id);
