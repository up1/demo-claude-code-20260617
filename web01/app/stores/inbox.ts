import axios, { type AxiosInstance } from 'axios'
import { defineStore } from 'pinia'

export type Channel = 'facebook' | 'line' | 'instagram'
export type Status = 'pending' | 'replied'

export interface InboxMessage {
  id: string
  customer_id: string
  sender_name: string
  avatar_url: string
  channel: Channel
  preview: string
  status: Status
  unread: boolean
  created_at: string
  updated_at: string
}

export interface Pagination {
  page: number
  page_size: number
  total_items: number
  total_pages: number
}

interface ListInboxResponse {
  success: boolean
  data: InboxMessage[]
  pagination: Pagination
}

// In a real app the JWT would come from an auth store / cookie. We keep a
// single source of truth here so the Authorization header is always attached.
const AUTH_TOKEN = 'demo-jwt-token'

let client: AxiosInstance | null = null

function getClient(): AxiosInstance {
  if (!client) {
    const config = useRuntimeConfig()
    client = axios.create({
      baseURL: config.public.apiBaseUrl || '',
      headers: { Authorization: `Bearer ${AUTH_TOKEN}` }
    })
  }
  return client
}

export const useInboxStore = defineStore('inbox', {
  state: () => ({
    messages: [] as InboxMessage[],
    channel: '' as '' | Channel,
    status: '' as '' | Status,
    search: '',
    page: 1,
    pageSize: 20,
    totalItems: 0,
    totalPages: 0,
    loading: false,
    unauthorized: false,
    error: ''
  }),
  getters: {
    isEmpty(state): boolean {
      return !state.loading && !state.error && state.messages.length === 0
    }
  },
  actions: {
    async fetchMessages() {
      this.loading = true
      this.unauthorized = false
      this.error = ''

      const params: Record<string, string | number> = {
        page: this.page,
        page_size: this.pageSize
      }
      if (this.channel) params.channel = this.channel
      if (this.status) params.status = this.status
      if (this.search) params.q = this.search

      try {
        const { data } = await getClient().get<ListInboxResponse>(
          '/api/v1/inbox/messages',
          { params }
        )
        this.messages = data.data
        this.totalItems = data.pagination.total_items
        this.totalPages = data.pagination.total_pages
        this.page = data.pagination.page
        this.pageSize = data.pagination.page_size
      } catch (err) {
        this.messages = []
        this.totalItems = 0
        this.totalPages = 0
        if (axios.isAxiosError(err) && err.response?.status === 401) {
          this.unauthorized = true
          this.error = 'Your session has expired. Please sign in again.'
        } else {
          this.error = 'Failed to load messages. Please try again.'
        }
      } finally {
        this.loading = false
      }
    },

    setChannel(channel: '' | Channel) {
      this.channel = channel
      this.page = 1
      return this.fetchMessages()
    },

    setStatus(status: '' | Status) {
      this.status = status
      this.page = 1
      return this.fetchMessages()
    },

    setSearch(search: string) {
      this.search = search
      this.page = 1
      return this.fetchMessages()
    },

    goToPage(page: number) {
      if (page < 1 || (this.totalPages && page > this.totalPages)) return
      this.page = page
      return this.fetchMessages()
    }
  }
})
