Feature: Draft Post Filtering
  As a content publisher
  I want draft posts to be skipped during site generation
  So that unpublished posts do not appear on the live site

  Scenario: Exclude draft posts from build outputs and sitemap
    Given a configuration directory with a valid profile
    And a blog directory containing 1 published post and 1 draft post
    When the build pipeline is executed
    Then the output directory should contain "blog/published-1.html"
    And the output directory should not contain "blog/draft-1.html"
    And the output file "sitemap.xml" should contain "blog/published-1.html"
    And the output file "sitemap.xml" should not contain "blog/draft-1.html"

  Scenario: Exclude draft posts from the generated RSS feed
    Given a configuration directory with a valid profile
    And a blog directory containing 1 published post and 1 draft post
    When the build pipeline is executed
    Then the output file "rss.xml" should contain "Published Post 1"
    And the output file "rss.xml" should not contain "Draft Post 1"

  Scenario: Exclude draft posts from search index and manifest registries
    Given a configuration directory with a valid profile
    And a blog directory containing 1 published post and 1 draft post
    When the build pipeline is executed
    Then the output file "search-index.json" should contain "Published Post 1"
    And the output file "search-index.json" should not contain "Draft Post 1"
    And the output file "api/manifest.json" should contain "Published Post 1"
    And the output file "api/manifest.json" should not contain "Draft Post 1"
