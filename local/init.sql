CREATE TABLE messages (
                          id SERIAL PRIMARY KEY,
                          phone VARCHAR(30) NOT NULL,
                          content TEXT NOT NULL CHECK (char_length(content) <= 100),
                          status VARCHAR(20) NOT NULL DEFAULT 'unsent',
                          created_at TIMESTAMP DEFAULT now(),
                          updated_at TIMESTAMP DEFAULT now(),
                          sent_at TIMESTAMP NULL,
                          remote_message_id TEXT NULL
);


INSERT INTO messages (phone, content)
VALUES
    ('+905301112233', 'Hello, this is the first test message.'),
    ('+905301112234', '2. Test message, for automatic sending.'),
    ('+905301112235', '3. Test message, for automatic sending.'),
    ('+905301112236', '4. Test message, for automatic sending.'),
    ('+905301112237', '5. Test message, for automatic sending.');