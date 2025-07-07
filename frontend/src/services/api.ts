import axios from 'axios';
import type {
  User,
  Activity,
  Group,
  Comment,
  LoginRequest,
  LoginResponse,
  UserCreateRequest,
  ActivityCreateRequest,
  GroupCreateRequest,
  CommentCreateRequest,
  CommentsResponse,
  GroupMembersResponse,
  GroupInvite,
} from '../types/api';

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

class ApiService {
  private api;

  constructor() {
    this.api = axios.create({
      baseURL: API_BASE_URL,
      headers: {
        'Content-Type': 'application/json',
      },
    });

    // Add token to requests if available
    this.api.interceptors.request.use((config) => {
      const token = localStorage.getItem('token');
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }
      return config;
    });
  }

  // Auth endpoints
  async login(credentials: LoginRequest): Promise<LoginResponse> {
    const response = await this.api.post<LoginResponse>('/login', credentials);
    return response.data;
  }

  async register(userData: UserCreateRequest): Promise<User> {
    const response = await this.api.post<User>('/users', userData);
    return response.data;
  }

  // User endpoints
  async getUser(userId: string): Promise<User> {
    const response = await this.api.get<User>(`/users/${userId}`);
    return response.data;
  }

  async getUserActivities(userId: string): Promise<Activity[]> {
    const response = await this.api.get<Activity[]>(`/users/${userId}/activities`);
    return response.data;
  }

  async getUserGroups(userId: string): Promise<Group[]> {
    const response = await this.api.get<Group[]>(`/users/${userId}/groups`);
    return response.data;
  }

  // Activity endpoints
  async createActivity(activityData: ActivityCreateRequest): Promise<Activity> {
    const response = await this.api.post<Activity>('/activities', activityData);
    return response.data;
  }

  async getUserFeed(userId: string): Promise<Activity[]> {
    const response = await this.api.get<Activity[]>(`/activities/feed?user_id=${userId}`);
    return response.data;
  }

  async getActivity(activityId: string): Promise<Activity> {
    const response = await this.api.get<Activity>(`/activities/${activityId}`);
    return response.data;
  }

  async updateActivity(activityId: string, activityData: Partial<ActivityCreateRequest>): Promise<Activity> {
    const response = await this.api.put<Activity>(`/activities/${activityId}`, activityData);
    return response.data;
  }

  async deleteActivity(activityId: string, creatorId: string): Promise<void> {
    await this.api.delete(`/activities/${activityId}`, {
      data: { creator_id: creatorId }
    });
  }

  // Group endpoints
  async createGroup(groupData: GroupCreateRequest): Promise<Group> {
    const response = await this.api.post<Group>('/groups', groupData);
    return response.data;
  }

  async getGroup(groupId: string, requesterId: string): Promise<Group> {
    const response = await this.api.get<Group>(`/groups/${groupId}?requester_id=${requesterId}`);
    return response.data;
  }

  async updateGroup(groupId: string, groupData: Partial<GroupCreateRequest>): Promise<Group> {
    const response = await this.api.put<Group>(`/groups/${groupId}`, groupData);
    return response.data;
  }

  async deleteGroup(groupId: string, creatorId: string): Promise<void> {
    await this.api.delete(`/groups/${groupId}`, {
      data: { creator_id: creatorId }
    });
  }

  async getGroupActivities(groupId: string, requesterId: string): Promise<Activity[]> {
    const response = await this.api.get<Activity[]>(`/groups/${groupId}/activities?requester_id=${requesterId}`);
    return response.data;
  }

  async getGroupMembers(groupId: string, requesterId: string): Promise<GroupMembersResponse> {
    const response = await this.api.get<GroupMembersResponse>(`/groups/${groupId}/members?requester_id=${requesterId}`);
    return response.data;
  }

  async addUserToGroup(groupId: string, userId: string): Promise<void> {
    await this.api.post(`/groups/${groupId}/members`, { user_id: userId });
  }

  async removeUserFromGroup(groupId: string, userId: string, requesterId: string): Promise<void> {
    await this.api.delete(`/groups/${groupId}/members`, {
      data: { user_id: userId, requester_id: requesterId }
    });
  }

  async createGroupInvite(groupId: string, creatorId: string, expiresAt: string): Promise<GroupInvite> {
    const response = await this.api.post<GroupInvite>(`/groups/${groupId}/invites`, {
      creator_id: creatorId,
      expires_at: expiresAt
    });
    return response.data;
  }

  async getGroupInvites(groupId: string): Promise<GroupInvite[]> {
    const response = await this.api.get<GroupInvite[]>(`/groups/${groupId}/invites`);
    return response.data;
  }

  async joinGroupByInvite(inviteCode: string, userId: string): Promise<void> {
    await this.api.post(`/invites/${inviteCode}/join`, { user_id: userId });
  }

  // Comment endpoints
  async getActivityComments(activityId: string): Promise<CommentsResponse> {
    const response = await this.api.get<CommentsResponse>(`/activities/${activityId}/comments`);
    return response.data;
  }

  async createComment(activityId: string, commentData: CommentCreateRequest): Promise<Comment> {
    const response = await this.api.post<Comment>(`/activities/${activityId}/comments`, commentData);
    return response.data;
  }

  async deleteComment(commentId: string, requesterId: string): Promise<void> {
    await this.api.delete(`/comments/${commentId}`, {
      data: { requester_id: requesterId }
    });
  }
}

export const apiService = new ApiService();
