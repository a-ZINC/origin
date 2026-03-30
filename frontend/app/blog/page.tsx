import { getPosts } from "@/lib/api";
import type { Post } from "@/lib/api";
import Link from "next/link";
import Tag  from "@/components/tag";

function fmt(d: string | null) {
  if (!d) return "";
  return new Date(d).toLocaleDateString("en-US", {
    month: "short", day: "numeric", year: "numeric",
  });
}

function PostCard({ post }: { post: Post }) {
  return (
    <Link href={`/blog/${post.Slug}`}
          className="block bg-s1 border border-b1 rounded-2xl p-7
                     hover:border-b2 hover:bg-s2 transition-all
                     duration-300 hover:-translate-y-0.5 group">

      {post.SeriesName && (
        <div className="inline-flex items-center gap-1.5 text-[11px] text-p2
                        bg-[rgba(139,92,246,0.1)] border border-[rgba(139,92,246,0.2)]
                        px-3 py-1 rounded-full font-mono mb-3">
          {post.SeriesName}
          {post.SeriesPos && ` · part ${post.SeriesPos}`}
        </div>
      )}

      <div className="flex gap-1.5 mb-4 flex-wrap">
        {post.Tags?.map((t) => <Tag key={t} name={t} />)}
        {post.Visibility === "private" && (
          <span className="text-[11px] px-2.5 py-0.5 rounded-[6px]
                           bg-[rgba(255,255,255,0.05)] text-t3 font-mono
                           border border-b1">
            private
          </span>
        )}
      </div>

      <h2 className="font-serif text-[21px] leading-[1.25] tracking-[-0.2px]
                     text-t1 mb-3 group-hover:text-p3 transition-colors">
        {post.Title}
      </h2>

      {post.Excerpt && (
        <p className="text-[13px] text-t2 leading-[1.7] mb-6 line-clamp-2">
          {post.Excerpt}
        </p>
      )}

      <div className="flex items-center gap-4 text-[12px] text-t3 font-mono
                      pt-5 border-t border-b1">
        <span>{fmt(post.PublishedAt)}</span>
        <span className="flex items-center gap-1">
          <svg className="w-3 h-3" viewBox="0 0 12 12" fill="none">
            <circle cx="6" cy="6" r="5" stroke="currentColor" strokeWidth="1"/>
            <path d="M6 3v3l2 1.5" stroke="currentColor" strokeWidth="1"
                  strokeLinecap="round"/>
          </svg>
          {post.ReadingTime} min
        </span>
        <span className="ml-auto text-p2 opacity-0
                         group-hover:opacity-100 transition-opacity">
          read →
        </span>
      </div>
    </Link>
  );
}

function TagFilter({ posts, activeTag }: { posts: Post[]; activeTag?: string }) {
  const tags = Array.from(new Set(posts.flatMap((p) => p.Tags ?? [])));
  if (!tags.length) return null;

  return (
    <div className="flex gap-2 flex-wrap mb-10">
      <a href="/blog"
         className={`text-[12px] px-4 py-1.5 rounded-full border
                     transition-all font-mono ${
           !activeTag
             ? "bg-[rgba(139,92,246,0.15)] border-[rgba(139,92,246,0.3)] text-p2"
             : "border-b2 text-t3 hover:text-t2 hover:border-b3"
         }`}>
        all
      </a>
      {tags.map((t) => (
        <a key={t} href={`/blog?tag=${t}`}
           className={`text-[12px] px-4 py-1.5 rounded-full border
                       transition-all font-mono ${
             activeTag === t
               ? "bg-[rgba(139,92,246,0.15)] border-[rgba(139,92,246,0.3)] text-p2"
               : "border-b2 text-t3 hover:text-t2 hover:border-b3"
           }`}>
          {t}
        </a>
      ))}
    </div>
  );
}

export default async function BlogPage({
  searchParams,
}: {
  searchParams: { tag?: string };
}) {
  const tag   = searchParams.tag;
  const posts = await getPosts(tag).catch(() => []);

  return (
    <div className="py-16">

      <div className="mb-12">
        <div className="text-[11px] tracking-[2.5px] uppercase text-t3
                        font-mono mb-3">writing</div>
        <h1 className="font-serif text-[48px] tracking-[-1px] text-t1 mb-3">
          {tag ? (
            <>All posts tagged <span className="text-p2 italic">#{tag}</span></>
          ) : (
            "All posts"
          )}
        </h1>
        <p className="text-[14px] text-t3 font-mono">
          {posts?.length ?? 0} posts{tag ? ` tagged "${tag}"` : ""}
        </p>
      </div>

      <TagFilter posts={posts ?? []} activeTag={tag} />

      {posts && posts.length > 0 ? (
        <div className="grid grid-cols-2 gap-4">
          {posts.map((p) => (
            <PostCard key={p.ID} post={p} />
          ))}
        </div>
      ) : (
        <div className="text-[14px] text-t3 font-mono py-20 text-center
                        border border-b1 rounded-2xl">
          {tag ? `no posts tagged "${tag}" yet` : "no posts yet"}
        </div>
      )}

    </div>
  );
}