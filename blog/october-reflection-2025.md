---
title: "October Reflection 2025"
description: "Reflections on evolving a Flask API from SQLite MVP to PostgreSQL with Docker, Azure cloud deployment, feature flags, and DevOps workflows. Read the full guide to learn."
date: 2025-10-28
tags: ["monthly-log", "python"]
---

## From SQLite to Azure: Practicing Production Thinking

I started this school management Flask API project with a simple goal: to learn how to build a REST API using SQLite as the MVP. It manages entities like students, instructors, and classes — making it a practical way to explore database relationships and the MVC architecture while testing each route.

As the project evolved, I realized it lacked a production-grade environment, which inspired me to expand it further:

- Migrating to PostgreSQL for a more realistic database setup.  
- Using feature flags via environment variables to enable gradual rollout.  
- Deploying to Azure for cloud reliability and scalability.
- Implementing CI/CD workflows with testing and coverage reporting.

What began as a simple learning exercise gradually became a chance to practice production thinking — designing, testing, and deploying like a real-world system.

---

## Starting with SQLite

The project began with SQLite, keeping the setup lightweight and ideal for an MVP. This allowed me to focus on fundamentals: designing routes, testing CRUD operations, and understanding how Flask handled requests and database interactions — without worrying about infrastructure.

This phase was about building confidence in the basics: getting the MVC pattern right, validating input, and shaping consistent responses. Running everything locally made it quick to iterate and see changes immediately.

Once the core API was stable, I began thinking about how a real system would run in production — and that’s when the project started to evolve.

---

## Migrating to PostgreSQL

After building a working MVP with SQLite, I wanted to see how the API would behave in a more realistic, production-style setup. PostgreSQL was the natural next step — widely used in industry, with features like connection pooling, schema management, and better concurrency handling.

I ran PostgreSQL locally using Docker, which made it easy to spin up a database without manual setup. This gave me hands-on experience with environment variables, connection strings, and configuration management in a real deployment scenario.

For production, I connected the API to Azure Database for PostgreSQL and deployed it to Azure Web App. This provided a reliable, persistent environment while introducing cloud concepts like managed databases, environment configuration, and service networking.

I also built a CI/CD pipeline that automatically checked code formatting, ran tests, and posted coverage reports in pull requests. These workflows brought a taste of DevOps practices, showing how automation supports reliability and feedback.

Finally, I explored pushing Docker images to Azure Container Registry (ACR), and discovered that GitHub provides its own container registry, offering an alternative workflow. This helped me connect the dots between registries, pipelines, and deployment targets.

### Imagining Gradual Rollout with Feature Flag

Drawing on my experience with feature flags from my internship at Shopify, I wanted to apply the same concepts to this personal project. Even without real users, I could simulate a gradual rollout by using an environment variable to toggle between SQLite and PostgreSQL.

In production, feature flags allow teams to enable changes incrementally, monitor behavior, and reduce risk. While I couldn’t fully replicate live traffic, this exercise let me practice the mindset and discipline behind safe rollouts — applying lessons learned in a professional environment to my own project.

### Removing the Feature Flag Gradually

Once all routes were fully migrated to PostgreSQL, I removed the environment variable that toggled between databases. This step mirrored a real-world cleanup process, where temporary feature flags are retired after a rollout is complete and stable.

Drawing on the practices I observed during my Shopify internship, I treated the removal as an intentional part of the lifecycle: plan the rollout, monitor behavior, and clean up temporary controls once they’re no longer needed. Feature flags are powerful tools, but leaving them in place too long can increase complexity and make the system harder to maintain.

Even in a personal project, thinking through the full lifecycle of temporary mechanisms — from introduction to removal — reinforced the discipline of maintaining clean, predictable, and maintainable systems.

---

## Lessons Learned

This project reinforced several important lessons about **applying production thinking** to even a personal, solo project:

- **Commit and PR granularity matters.** Small, focused changes are easier to review, test, and rollback if needed. While I initially grouped multiple commits in a single pull request, treating each route or feature as its own PR would mirror best practices in real engineering teams.  

- **Feature flags enable safe experimentation.** Leveraging the experience I gained at Shopify, I applied environment-based toggles to simulate gradual rollout and controlled migration. Thinking through their full lifecycle — from introduction to removal — reinforced disciplined, maintainable workflows.  

- **Solo projects can simulate production practices.** Even without a team or real users, you can practice DevOps principles: phased migrations, deployment safety, cloud integration, and incremental changes.  

- **Understanding the ecosystem is valuable.** Exploring container registries like Azure Container Registry and GitHub Container Registry helped me see how CI/CD pipelines, images, and deployments connect — a key piece of real-world operations.  

- **CI/CD accelerates feedback and quality.** Automated testing, code formatting checks, and coverage reports in pull requests mimic real engineering workflows, providing faster feedback and improving reliability.

---

## Thank you

Big thanks for reading! You’re awesome, and I hope this post helped. Until next time!
