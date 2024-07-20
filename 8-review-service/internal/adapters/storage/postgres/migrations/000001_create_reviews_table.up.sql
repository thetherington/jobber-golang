CREATE TABLE IF NOT EXISTS public.reviews (
  id SERIAL UNIQUE,
  gig_id text NOT NULL,
  reviewer_id text NOT NULL,
  order_id text NOT NULL,
  seller_id text NOT NULL,
  review text NOT NULL,
  reviewer_image text NOT NULL,
  reviewer_username text NOT NULL,
  country text NOT NULL,
  review_type text NOT NULL,
  rating integer DEFAULT 0 NOT NULL,
  created_at timestamp DEFAULT CURRENT_DATE,
  PRIMARY KEY (id)
);

CREATE INDEX IF NOT EXISTS gig_id_idx ON public.reviews (gig_id);

CREATE INDEX IF NOT EXISTS seller_id_idx ON public.reviews (seller_id);