# Kargo in Action — companion code

Example YAML and configuration from **_Kargo in Action: Continuous promotion for Kubernetes_** (Manning), by Ken Cochrane, Jesse Suen, and Kent Rancourt. Examples are organized by chapter and match the listings in the book.

> **MEAP note:** the book is in early access and under active development. Chapters 1–3 are the current review focus; later-chapter examples are still being finalized.

## Layout

| Folder | Chapter |
|---|---|
| [`ch03/`](ch03) | Ch 3 — Getting started with Kargo (install + first promotion) |
| [`ch04/`](ch04) | Ch 4 — Modeling your delivery: Warehouses, Freight, and Stage graphs |
| [`ch05/`](ch05) | Ch 5 — Promoting safely: automatic, manual, and recoverable |
| [`ch07/`](ch07) | Ch 7 — Scaling Kargo across services and teams |
| [`appendix-a/`](appendix-a) | Appendix A — Kargo resource reference |
| [`appendix-b/`](appendix-b) | Appendix B — Built-in promotion step reference |

Chapters 1, 2, 6, 8, 9, and 10 are conceptual or narrative and have no standalone code.

## Running the examples

The getting-started walkthrough in chapter 3 promotes the sample application **demo-app** through a GitOps repository that Argo CD reconciles. It uses these companion repositories:

- **[demo-app](https://github.com/kargobook/demo-app)** — the sample application image (`ghcr.io/kargobook/demo-app`)
- **[demo-gitops](https://github.com/kargobook/demo-gitops)** — the GitOps configuration repository Kargo writes to and Argo CD syncs

Later chapters also reference the **[payments-api](https://github.com/kargobook/payments-api)** and **[auth-service](https://github.com/kargobook/auth-service)** sample services for multi-service examples.

You will need a Kubernetes cluster with Argo CD and Kargo installed; chapter 3 covers the full installation.

## Feedback

Found a problem with an example? Please open an issue. Book feedback goes through the Manning liveBook forum.
