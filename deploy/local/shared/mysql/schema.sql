CREATE TABLE IF NOT EXISTS monster (
  id VARCHAR(36) NOT NULL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  health INT(11) NOT NULL,
  max_health INT(11) NOT NULL,
  attack INT(11) NOT NULL,
  defense INT(11) NOT NULL,
  speed INT(11) NOT NULL,
  avatar_url TEXT NOT NULL,
  is_partnerable TINYINT(1) NOT NULL,
  KEY `is_partnerable` (`is_partnerable`)
);

CREATE TABLE IF NOT EXISTS game (
  id VARCHAR(36) NOT NULL PRIMARY KEY,
  player_name VARCHAR(255) NOT NULL,
  created_at BIGINT(20) NOT NULL,
  battle_won INT(11) NOT NULL,
  scenario VARCHAR(30) NOT NULL,
  partner_id VARCHAR(36) NOT NULL
);

CREATE TABLE IF NOT EXISTS battle (
  game_id VARCHAR(36) NOT NULL PRIMARY KEY,
  state VARCHAR(30) NOT NULL,
  partner_monster_id VARCHAR(36) NOT NULL,
  partner_name VARCHAR(255) NOT NULL,
  partner_max_health INT(11) NOT NULL,
  partner_health INT(11) NOT NULL,
  partner_attack INT(11) NOT NULL,
  partner_defense INT(11) NOT NULL,
  partner_speed INT(11) NOT NULL,
  partner_avatar_url TEXT NOT NULL,
  partner_last_damage INT(11) NOT NULL,
  enemy_monster_id VARCHAR(36) NOT NULL,
  enemy_name VARCHAR(255) NOT NULL,
  enemy_max_health INT(11) NOT NULL,
  enemy_health INT(11) NOT NULL,
  enemy_attack INT(11) NOT NULL,
  enemy_defense INT(11) NOT NULL,
  enemy_speed INT(11) NOT NULL,
  enemy_avatar_url TEXT NOT NULL,
  enemy_last_damage INT(11) NOT NULL
);

CREATE TABLE IF NOT EXISTS user (
  id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  username VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL UNIQUE,
  password VARCHAR(255) NOT NULL,
  created_at BIGINT(20) NOT NULL
);

CREATE TABLE IF NOT EXISTS event (
  id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(255) NOT NULL
);

CREATE TABLE venue (
  id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  open_days VARCHAR(15) NOT NULL,  -- To store the days of the week as a comma-separated list
  open_at VARCHAR(5) NOT NULL,
  closed_at VARCHAR(5) NOT NULL,
  timezone VARCHAR(255) NOT NULL
);

CREATE TABLE venue_event (
  venue_id INT NOT NULL,
  event_id INT NOT NULL,
  meetups_capacity INT NOT NULL,
  PRIMARY KEY (venue_id, event_id),
  FOREIGN KEY (venue_id) REFERENCES venue(id),
  FOREIGN KEY (event_id) REFERENCES event(id)
);
