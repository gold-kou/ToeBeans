INSERT INTO users (name, email, password, activation_key, email_verified) VALUES('guest', 'guestUser@example.com' , 'Guest1234', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', true);

-- ハッシュ化したパスワードでないため、ログインする場合は/loginのパスワードチェック構文をコメントアウトする。
INSERT INTO users (name, email, password, activation_key, email_verified) VALUES('userA', 'userA@example.com' , 'userA1234', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', true);