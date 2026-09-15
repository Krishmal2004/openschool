<p align="center">
  <img src="frontend/public/w-favicon.webp" alt="OpenSchool" width="96" />
</p>

<h1 align="center">OpenSchool</h1>

<p align="center">
  Open digital platform for Sri Lankan schools.
</p>

<p align="center">
  <a href="https://github.com/openschool-org/openschool/actions/workflows/backend-ci.yml"><img src="https://github.com/openschool-org/openschool/actions/workflows/backend-ci.yml/badge.svg" alt="Backend CI"></a>
  <a href="https://github.com/openschool-org/openschool/actions/workflows/frontend-ci.yml"><img src="https://github.com/openschool-org/openschool/actions/workflows/frontend-ci.yml/badge.svg" alt="Frontend CI"></a>
  <a href="./LICENSE"><img src="https://img.shields.io/badge/license-Apache%202.0-blue.svg" alt="Apache 2.0 License"></a>
  <a href="./CODE_OF_CONDUCT.md"><img src="https://img.shields.io/badge/Contributor%20Covenant-2.1-4baaaa.svg" alt="Contributor Covenant"></a>
</p>

---

## About OpenSchool

OpenSchool is an open-source school management platform for Sri Lankan government schools. It brings everyday school work into one place: student and staff records, classes, attendance, curriculum, timetables, marks, notifications, reports, and school administration.

Many schools still depend on paper records or several disconnected tools. This makes information harder to find, update, and share safely. OpenSchool is built to give schools one clear system for managing that work.

OpenSchool can be self-hosted, so a school or public-sector organisation can run it on infrastructure it controls. Student and administrative data can remain under local control instead of depending on a commercial foreign-hosted service. It is licensed under Apache License 2.0, with no licence fees and no vendor lock-in.

## Features

- School setup, academic years, grades, classes, houses, and streams
- Student, teacher, guardian, and non-academic staff records
- Daily student and staff attendance
- Subjects, curriculum choices, term marks, and promotions
- Timetable planning, review, and publishing
- In-app notifications, reports, dashboards, and automation checks
- Separate portals for administrators, teachers, students, and parents

Read the [feature guide](docs/FEATURES.md) for the full list.

## Built with

| Part | Technology |
| --- | --- |
| Backend | Go and Gin |
| Frontend | React, TypeScript, Vite, and Carbon Design System |
| Database | PostgreSQL |
| Identity and access management | [ThunderID](https://github.com/thunderid) |

The backend is a modular monolith: one application, one database, and clear feature modules. See the [architecture guide](docs/ARCHITECTURE.md) for more detail.

## Get started

You need Go, Node.js with pnpm, Docker, and Docker Compose.

```bash
git clone https://github.com/openschool-org/openschool.git
cd openschool
```

Then follow these guides:

1. [Set up ThunderID](docs/THUNDERID.md).
2. [Set up and run OpenSchool](docs/SETUP.md).
3. Open `http://localhost:5173` and create the first administrator account.

## Documentation

| Guide | Use it for |
| --- | --- |
| [Setup](docs/SETUP.md) | Running OpenSchool locally and setting up a new school |
| [Features](docs/FEATURES.md) | Understanding what the system can do |
| [Architecture](docs/ARCHITECTURE.md) | Understanding the project structure and data model |
| [Architecture decisions](docs/adr/) | Understanding why key technical decisions were made |
| [Contributing](CONTRIBUTING.md) | Setting up a development environment and opening a pull request |

## Contributing

Contributions are welcome. Please read [CONTRIBUTING.md](CONTRIBUTING.md) before opening an issue or pull request, and follow the [Code of Conduct](CODE_OF_CONDUCT.md).

## Security

If you find a security issue, please do not report it in a public issue. Follow [SECURITY.md](SECURITY.md).

## License

OpenSchool is licensed under the [Apache License 2.0](LICENSE).
