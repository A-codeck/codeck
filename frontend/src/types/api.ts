// User types
export interface User {
  id: string;
  name: string;
  email: string;
}

// Activity types
export interface Activity {
  id: string;
  title: string;
  description: string;
  activity_image?: string;
  date: string;
  creator_id: string;
}

// Group types
export interface Group {
  id: string;
  name: string;
  description: string;
  group_image?: string;
  start_date?: string;
  end_date: string;
  creator_id: string;
}

export interface GroupMember {
  user_id: string;
  group_id: string;
  nickname?: string;
}

export interface GroupInvite {
  invite_code: string;
  group_id: string;
  created_by: string;
  created_at: string;
  expires_at: string;
  is_active: boolean;
}

// Comment types
export interface Comment {
  id: string;
  content: string;
  user_id: string;
  activity_id: string;
  created_at: string;
}

// Request types
export interface LoginRequest {
  email: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  user: User;
}

export interface UserCreateRequest {
  name: string;
  email: string;
  password: string;
}

export interface ActivityCreateRequest {
  title: string;
  description: string;
  activity_image?: string;
  date: string;
  group_id?: string;
  creator_id: string; // Required by backend
}

export interface GroupCreateRequest {
  name: string;
  description: string;
  group_image?: string;
  end_date: string;
  creator_id: string; // Required by backend
}

export interface CommentCreateRequest {
  content: string;
  user_id: string;
}

// Response types
export interface ErrorResponse {
  error: string;
}

export interface CommentsResponse {
  activity_id: string;
  comment_count: number;
  comments: Comment[];
}

export interface GroupMembersResponse {
  group_id: string;
  member_count: number;
  members: GroupMember[];
}

// Extended types for frontend
export interface ActivityWithGroup extends Activity {
  group?: Group;
  comments?: Comment[];
  commentsCount?: number;
}

export interface UserStats {
  user_id: string;
  user_name: string;
  activity_count: number;
}
