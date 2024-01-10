export interface Feed {
  id: string
  name: string
  type: string
  feedURL: string
  hostURL: string
  author: string
  description: string
  imageURL: string
  accountID: string
  account: {
    username: string
    uri: string
  }
}

export interface Status {
  createdAt: string
  content: string
}
