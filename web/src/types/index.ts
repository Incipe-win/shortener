// API types shared across the application

export interface ConvertRequest {
  long_url: string;
}

export interface ConvertResponse {
  short_url: string;
}

export interface PreviewResponse {
  short_url: string;
  long_url: string;
  summary: string;
  keywords: string[];
  risk_level: string;
  risk_reason?: string;
}

export interface LinkItem {
  id: number;
  surl: string;
  lurl: string;
  ai_summary: string;
  ai_keywords: string[];
  risk_level: string;
  click_count: number;
  create_at: string;
}

export interface LinksResponse {
  list: LinkItem[];
  total: number;
  page: number;
  page_size: number;
}
