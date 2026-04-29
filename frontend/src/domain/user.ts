export type Role = 'user' | 'admin';

export interface UserProfile {
  id: string;
  email: string;
  roles: Role[];
  verified: boolean;
  createdAt: string;
}
