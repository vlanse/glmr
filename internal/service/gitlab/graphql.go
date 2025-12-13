package gitlab

const (
	projectMRsRequest = `
{
  project(fullPath: "%s") {
    mergeRequests(state: opened) {
      nodes {
        id
        iid
        projectId
        createdAt
        updatedAt
        webUrl
        conflicts
		draft
        title
        state
        committers {
          nodes {
            id
            username
            avatarUrl
            webUrl
            name
            publicEmail
          }
        }
        approvedBy {
          nodes {
            id
            username
            avatarUrl
            webUrl
            name
            publicEmail
          }
        }
        author {
          id
          username
          avatarUrl
          webUrl
          name
          publicEmail
        }
        headPipeline {
          id
          status
          path
        }
        diffStatsSummary {
          additions
          changes
          fileCount
          deletions
        }
      }
    }
  }
}
`
)
