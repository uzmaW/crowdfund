export type Project = {
  id: string;
  title: string;
  description: string;
  fundingGoal: number;
  currentFunding: number;
  deadline: Date;
  status: 'active' | 'completed' | 'canceled';
  creatorId: string;
  createdAt: Date;
  updatedAt: Date;
}