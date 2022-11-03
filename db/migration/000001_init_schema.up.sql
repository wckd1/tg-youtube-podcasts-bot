CREATE TABLE "subscriptions" (
  "id" INTEGER PRIMARY KEY,
  "source_path" TEXT NOT NULL,
  "source_type" TEXT NOT NULL,
  "title" TEXT,
  UNIQUE("source_path","title")
);
