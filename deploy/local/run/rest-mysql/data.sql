INSERT INTO monster (id, name, max_health, health, attack, defense, speed, avatar_url, is_partnerable) VALUES
  ("b1c87c5c-2ac3-471d-9880-4812552ee15d", 'Yellowleg', 100, 100, 25, 5, 15, "https://haraj-sol-dev.s3.eu-west-1.amazonaws.com/hex-monscape/monsters/yellowleg.png", 1),
  ("0f9b84b6-a768-4ba9-8800-207740fc993d", 'Bluebub', 100, 100, 20, 15, 10, "https://haraj-sol-dev.s3.eu-west-1.amazonaws.com/hex-monscape/monsters/bluebub.png", 1),
  ("85db0102-212d-4ac8-932c-a0e876a29a85", 'Grumpy', 100, 100, 25, 5, 20, "https://haraj-sol-dev.s3.eu-west-1.amazonaws.com/hex-monscape/monsters/grumpy.png", 1),
  ("5e1ab413-415a-4326-8e39-0f56f8a66054", 'Vegiewee', 150, 150, 25, 20, 12, "https://haraj-sol-dev.s3.eu-west-1.amazonaws.com/hex-monscape/monsters/vegiewee.png", 0),
  ("c2ca1953-2376-489e-8e34-8bb48957f140", 'Snekworm', 150, 150, 30, 5, 21, "https://haraj-sol-dev.s3.eu-west-1.amazonaws.com/hex-monscape/monsters/snekworm.png", 0),
  ("88a98dee-ce84-4afb-b5a8-7cc07535f73f", 'Waneye', 100, 100, 20, 10, 15, "https://haraj-sol-dev.s3.eu-west-1.amazonaws.com/hex-monscape/monsters/waneye.png", 0);

INSERT INTO user (username, email, password, created_at) VALUES
  ("marion", "marion@eveners.com", "123456", UNIX_TIMESTAMP()),
  ("todd", "todd@eveners.com", "123456", UNIX_TIMESTAMP()),
  ("anthony", "anthony@eveners.com", "123456", UNIX_TIMESTAMP()),
  ("kevin", "kevin@eveners.com", "123456", UNIX_TIMESTAMP()),
  ("eric", "eric@eveners.com", "123456", UNIX_TIMESTAMP()),
  ("elnora", "elnora@eveners.com", "123456", UNIX_TIMESTAMP()),
  ("etta", "etta@eveners.com", "123456", UNIX_TIMESTAMP()),
  ("caleb", "caleb@eveners.com", "123456", UNIX_TIMESTAMP()),
  ("larry", "larry@eveners.com", "123456", UNIX_TIMESTAMP()),
  ("stanley", "stanley@eveners.com", "123456", UNIX_TIMESTAMP()),
  ("nelle", "nelle@eveners.com", "123456", UNIX_TIMESTAMP()),
  ("luke", "luke@eveners.com", "123456", UNIX_TIMESTAMP()),
  ("ian", "ian@eveners.com", "123456", UNIX_TIMESTAMP()),
  ("harry", "harry@eveners.com", "123456", UNIX_TIMESTAMP()),
  ("paul", "paul@eveners.com", "123456", UNIX_TIMESTAMP()),
  ("lula", "lula@eveners.com", "123456", UNIX_TIMESTAMP()),
  ("warren", "warren@eveners.com", "123456", UNIX_TIMESTAMP()),
  ("marguerite", "marguerite@eveners.com", "123456", UNIX_TIMESTAMP()),
  ("mitchell", "mitchell@eveners.com", "123456", UNIX_TIMESTAMP());

INSERT INTO event (name) VALUES
  ('Wedding'),
  ('Exhibition'),
  ('Bazaar'),
  ('Workshop'),
  ('Conference');
