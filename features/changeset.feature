Feature: Changesets

Background:
   Given I'm in project dir

Scenario: No changes
  Given file ./Oyafile containing
    """
    Project: project
    Changeset: echo ""
    all: |
      echo "Root"
    """
  And file ./project1/Oyafile containing
    """
    Changeset: echo ""
    all: |
      echo "Project1"
    """
  When I run "oya run --changeset --recurse all"
  Then the command succeeds
  And the command outputs
  """
  """

Scenario: Child marks itself as changed
  Given file ./Oyafile containing
    """
    Project: project
    Changeset: echo ""
    all: |
      echo "Root"
    """
  And file ./project1/Oyafile containing
    """
    Changeset: echo "+."
    all: |
      echo "Root"
    """
  When I run "oya run --changeset --recurse all"
  Then the command succeeds
  And the command outputs
  """
  Root

  """

Scenario: Child marks parent as changed
  Given file ./Oyafile containing
    """
    Project: project
    Changeset: echo ""
    all: |
      echo "Root"
    """
  And file ./project1/Oyafile containing
    """
    Changeset: echo "+../"
    all: |
      echo "Root"
    """
  When I run "oya run --changeset --recurse all"
  Then the command succeeds
  And the command outputs
  """
  Root

  """

Scenario: Parent marks child as changed
  Given file ./Oyafile containing
    """
    Project: project
    Changeset: echo "+project1/"
    all: |
      echo "Root"
    """
  And file ./project1/Oyafile containing
    """
    Changeset: echo ""
    all: |
      echo "Project1"
    """
  When I run "oya run --changeset --recurse all"
  Then the command succeeds
  And the command outputs
  """
  Project1

  """
