import type { Feed, Status } from '../models/models'
export async function fetchFeeds() {
  try {
    const response = await fetch('/api/feeds')
    if (!response.ok) {
      throw new Error('Failed to fetch data')
    }
    const jsonData = await response.json()
    return jsonData
  } catch (error) {
    console.error('Error fetching data:', error)
    return []
  }
}

export async function addFeed(feedURL: string): Promise<Feed | null> {
  try {
    const response = await fetch('/api/feeds', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        feedURL: feedURL
      })
    })

    if (!response.ok) {
      if (response.status === 400) {
        const jsonData = await response.json()
        throw new Error(jsonData.errors[0]?.message)
      }
      throw new Error('Failed to submit data')
    }

    const jsonData = await response.json()
    const feed: Feed = jsonData

    return feed
  } catch (error) {
    console.error('Error submitting data:', error)
    return null
  }
}

export async function fetchUserFeed(username: string): Promise<Feed | null> {
  try {
    const response = await fetch(`/api/users/${username}/feed`)
    if (!response.ok) {
      throw new Error('Failed to fetch user feed data')
    }
    const jsonData = await response.json()
    return jsonData
  } catch (error) {
    console.error('Error fetching user feed data:', error)
    return null
  }
}

export async function fetchFeedStatus(feedId: string | null): Promise<Status[] | null> {
  try {
    if (!feedId) {
      throw new Error('Feed ID is missing')
    }

    const response = await fetch(`/api/feeds/${feedId}/status?limit=10`)
    if (!response.ok) {
      throw new Error('Failed to fetch feed status data')
    }
    const jsonData = await response.json()
    return jsonData.items
  } catch (error) {
    console.error('Error fetching feed status data:', error)
    return null
  }
}
