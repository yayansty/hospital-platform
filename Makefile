up:
	docker compose up -d

down:
	docker compose down

restart:
	docker compose restart

build:
	docker compose build

ps:
	docker compose ps

logs:
	docker compose logs -f

shell:
	docker compose exec app bash

composer:
	docker compose exec app composer

artisan:
	docker compose exec app php artisan
