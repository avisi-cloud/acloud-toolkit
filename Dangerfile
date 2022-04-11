
if git.lines_of_code > 2_000
  warn "This merge request is definitely too big (more than #{git.lines_of_code} lines changed), please split it into multiple merge requests."
elsif git.lines_of_code > 500
  warn "This merge request is quite big (more than #{git.lines_of_code} lines changed), please consider splitting it into multiple merge requests."
end

if gitlab.mr_body.size < 5
    fail "Please provide a proper merge request description."
end

unless gitlab.mr_json["assignee"]
    warn "This merge request does not have any assignee yet. Setting an assignee clarifies who needs to take action on the merge request at any given time."
end

if gitlab.mr_author != "automation-user"
 
jira.check(
  key: ["[a-zA-Z]+"],
  url: "https://insight.avisi.nl/jira/browse",
  search_title: true,
  search_commits: true,
  fail_on_warning: true,
  report_missing: true,
  skippable: false
)
end

markdown("**If needed, you can retry the [`danger-review` job](#{ENV['CI_JOB_URL']}) that generated this comment.**")
