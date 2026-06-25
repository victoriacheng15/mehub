Feature: Site Generation Pipeline
  As a developer
  I want the generator to build the website
  So that my changes are compiled into static files

  Scenario: Build the site successfully
    Given a configuration directory with a valid profile
    And a blog directory containing 1 published post
    When the build pipeline is executed
    Then the output directory should contain "index.html"
    And the output directory should contain "api/manifest.json"

  Scenario: Build fails when configuration is missing
    Given a configuration directory with a missing profile
    When the build pipeline is executed
    Then the build pipeline execution should fail
