Feature: Site Generation Pipeline
  As a developer
  I want the generator to build the website
  So that my changes are compiled into static files

  Scenario: Build the site successfully
    Given a configuration directory with a valid profile
    And a blog directory containing 1 published post
    When the build pipeline is executed
    Then the output directory should contain "index.html"
    And the output directory should contain "about.html"
    And the output directory should contain "blog.html"
    And the output directory should contain "404.html"
    And the output directory should contain "sitemap.xml"
    And the output directory should contain "rss.xml"
    And the output directory should contain "search-index.json"
    And the output directory should contain "llms.txt"
    And the output directory should contain "api/manifest.json"
    And the output directory should contain "blog"
    And the output directory should contain "tags"

  Scenario: Copy static assets successfully
    Given a configuration directory with a valid profile
    And a static assets directory containing a file "robots.txt"
    And a blog directory containing 1 published post
    When the build pipeline is executed
    Then the output directory should contain "robots.txt"

  Scenario: Build fails when configuration is missing
    Given a configuration directory with a missing profile
    When the build pipeline is executed
    Then the build pipeline execution should fail

  Scenario: Render open source contributions on the home page
    Given a configuration directory with a valid profile containing contributions
    When the build pipeline is executed
    Then the output file "index.html" should contain "Open Source Contributions"
    And the output file "index.html" should contain "Last updated: 2026-07-18"
    And the output file "index.html" should contain "test-repo"
