import { getPost } from "@/lib/api";
import { notFound } from "next/navigation";
import Link from "next/link";
import Tag  from "@/components/tag";
import ViewTracker from "./view-tracker";

function fmt(d: string | null) {
  if (!d) return "";
  return new Date(d).toLocaleDateString("en-US", {
    month: "long", day: "numeric", year: "numeric",
  });
}

export async function generateMetadata({ params }: { params: Promise<{ slug: string }>;}) {
  try {
    const { slug } = await params;
    const post = await getPost(slug);
    return { title: post.Title, description: post.Excerpt };
  } catch {
    return { title: "Post not found" };
  }
}

export default async function PostPage({ params }: {params: Promise<{ slug: string }>;} ) {
  let post;
  try {
    const { slug } = await params;
    post = await getPost(slug);
  } catch {
    notFound();
  }

  return (
    <div className="py-12">

      <ViewTracker postId={post.ID} />

      <Link href="/blog"
            className="inline-flex items-center gap-1.5 text-[11px] text-txt-3
                       hover:text-txt-2 transition-colors mb-8 tracking-wide">
        ← all posts
      </Link>

      {post.SeriesName && (
        <Link href={`/series/${post.SeriesSlug}`}
              className="inline-flex items-center gap-2 text-[11px] text-v-2
                         bg-v/8 border border-v/18 rounded-[6px] px-3 py-1.5
                         mb-5 hover:bg-v/12 transition-colors">
          {post.SeriesName}
          {post.SeriesPos && (
            <span className="text-txt-3">· part {post.SeriesPos}</span>
          )}
        </Link>
      )}

      {/* title */}
      <h1 className="font-serif text-[36px] sm:text-[42px] font-normal
                     text-txt leading-[1.15] tracking-[-0.3px] mb-5">
        {post.Title}
      </h1>

      {/* meta bar */}
      <div className="flex items-center gap-4 py-4 border-t border-b
                      border-line mb-12 flex-wrap">
        <span className="flex items-center gap-1.5 text-[11px] text-txt-3">
          <svg className="w-3 h-3" viewBox="0 0 12 12" fill="none">
            <circle cx="6" cy="6" r="5" stroke="currentColor" strokeWidth="1"/>
            <path d="M6 3v3l2 1.5" stroke="currentColor" strokeWidth="1"
                  strokeLinecap="round"/>
          </svg>
          {post.ReadingTime} min read
        </span>

        <div className="w-px h-3 bg-line-2" />

        <span className="text-[11px] text-txt-3">{fmt(post.PublishedAt)}</span>

        <div className="w-px h-3 bg-line-2" />

        <div className="flex gap-1.5">
          {post.Tags?.map((t) => <Tag key={t} name={t} />)}
        </div>
      </div>

      {/* body — rendered HTML from Go */}
      <div
        className="
          prose prose-invert max-w-none
          prose-p:font-serif prose-p:text-[17px] prose-p:leading-[1.95]
          prose-p:text-[#a09dbd] prose-p:mb-7
          first-letter:text-[52px] first-letter:float-left
          first-letter:leading-[0.82] first-letter:mr-3.5
          first-letter:mt-1.5 first-letter:font-serif first-letter:text-txt
          prose-h2:font-serif prose-h2:text-[24px] prose-h2:font-normal
          prose-h2:text-txt prose-h2:mt-10 prose-h2:mb-4
          prose-h3:text-[11px] prose-h3:tracking-[2px] prose-h3:uppercase
          prose-h3:text-txt-3 prose-h3:font-mono prose-h3:font-normal
          prose-h3:mt-8 prose-h3:mb-3
          prose-code:text-v-3 prose-code:bg-v/8 prose-code:px-2
          prose-code:py-0.5 prose-code:rounded prose-code:text-[13px]
          prose-code:border prose-code:border-v/15 prose-code:font-mono
          prose-pre:bg-[#07070f] prose-pre:border prose-pre:border-line-2
          prose-pre:rounded-[10px] prose-pre:p-6
          prose-blockquote:border-l-2 prose-blockquote:border-v/40
          prose-blockquote:pl-5 prose-blockquote:text-txt-3
          prose-blockquote:italic prose-blockquote:font-serif
          prose-a:text-v-2 prose-a:no-underline hover:prose-a:text-v-3
          prose-strong:text-txt prose-strong:font-medium
          prose-li:font-serif prose-li:text-[16px] prose-li:text-[#a09dbd]
        "
        dangerouslySetInnerHTML={{ __html: post.HTMLBody }}
      />

      {/* series prev/next */}
      {(post.PrevPost || post.NextPost) && (
        <div className="grid grid-cols-2 gap-3 mt-14 pt-8 border-t border-line">
          {post.PrevPost ? (
            <Link href={`/blog/${post.PrevPost.Slug}`}
                  className="bg-bg-1 border border-line rounded-[10px] p-4
                             hover:border-v/30 hover:bg-bg-2 transition-all group">
              <div className="text-[9px] tracking-[1.5px] uppercase text-txt-3 mb-2">
                ← previous
              </div>
              <div className="text-[12px] text-txt-2 group-hover:text-txt transition-colors
                              leading-snug">
                {post.PrevPost.Title}
              </div>
            </Link>
          ) : <div />}

          {post.NextPost ? (
            <Link href={`/blog/${post.NextPost.Slug}`}
                  className="bg-bg-1 border border-line rounded-[10px] p-4
                             hover:border-v/30 hover:bg-bg-2 transition-all
                             group text-right">
              <div className="text-[9px] tracking-[1.5px] uppercase text-txt-3 mb-2">
                next →
              </div>
              <div className="text-[12px] text-txt-2 group-hover:text-txt transition-colors
                              leading-snug">
                {post.NextPost.Title}
              </div>
            </Link>
          ) : <div />}
        </div>
      )}

    </div>
  );
}