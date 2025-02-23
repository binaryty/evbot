CREATE TABLE IF NOT EXISTS events (
                                      id INTEGER PRIMARY KEY AUTOINCREMENT,
                                      user_id INTEGER NOT NULL,
                                      title TEXT NOT NULL,
                                      description TEXT,
                                      date DATETIME NOT NULL,
                                      created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS users (
                                     user_id INTEGER PRIMARY KEY,
                                     first_name TEXT NOT NULL,
                                     username TEXT,
                                     created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                                     updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS registrations (
                                             event_id INTEGER NOT NULL,
                                             user_id INTEGER NOT NULL,
                                             created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                                             PRIMARY KEY (event_id, user_id),
                                             FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS user_states (
                                           user_id INTEGER PRIMARY KEY,
                                           state_data TEXT NOT NULL,
                                           created_at DATETIME NOT NULL
);