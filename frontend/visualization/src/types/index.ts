export interface Person {
  member_id: number
  name: string
  style_name?: string
  gender: string
  generation: number
  lineage_path: string
  father_id?: number
  birth_year?: number
  death_year?: number
  birth_place?: string
  burial_place?: string
  biography?: string
  avatar_url?: string
  created_at?: string
  updated_at?: string
}

export interface GenealogyNode {
  id: number
  name: string
  generation: number
  gender: string
  children?: GenealogyNode[]
  spouse?: GenealogyNode[]
}

export interface Statistics {
  total_persons: number
  total_generations: number
  male_count: number
  female_count: number
  generation_distribution: { generation: number; count: number }[]
}
