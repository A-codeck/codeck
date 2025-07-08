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
  group_ids?: string[]; // Now supports multiple groups
}

// Activity with group information for feed display
export interface ActivityWithGroups {
  id: string;
  title: string;
  description: string;
  activity_image?: string;
  date: string;
  creator_id: string;
  group_names?: string[]; // Group names for display
}

// Group types
export interface Group {
  id: string;
  name: string;
  description: string;
  group_image?: string;
  creator_id: string;
}

export interface GroupMember {
  user_id: string;
  group_id: string;
  nickname?: string;
}

export interface GroupMemberWithUser {
  user_id: string;
  group_id: string;
  nickname?: string;
  user_name: string;
  user_email: string;
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
  user_name?: string;
  user_email?: string;
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
  date: string;
  group_ids: string[]; // Now required and supports multiple groups
  creator_id: string; // Required by backend
  image: File; // Now required - File object for upload
}

export interface GroupCreateRequest {
  name: string;
  description: string;
  creator_id: string; // Required by backend
  image?: File; // Now optional - File object for upload
}

export interface CommentCreateRequest {
  content: string;
  user_id: string;
}

export interface AddUserByEmailRequest {
  email: string;
  requester_id: string;
}

export interface LeaveGroupRequest {
  user_id: string;
}

export interface JoinGroupByInviteRequest {
  invite_code: string;
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

// Enhanced comments response with user information
export interface CommentsWithUsersResponse {
  activity_id: string;
  comment_count: number;
  comments: Comment[];
}

export interface GroupMembersResponse {
  group_id: string;
  member_count: number;
  members: GroupMemberWithUser[];
}

// Extended types for frontend
export interface ActivityWithGroup extends Activity {
  group?: Group;
  comments?: Comment[];
  commentsCount?: number;
}

// Extended types for activities with group information
export interface ActivityWithGroupNames extends ActivityWithGroups {
  group?: Group; // For backward compatibility
  comments?: Comment[];
  commentsCount?: number;
}

export interface UserStats {
  user_id: string;
  user_name: string;
  activity_count: number;
}
