import Link   from "next/link";
import Tag    from "@/components/tag";
import { getPosts, getProjects } from "@/lib/api";

const stack = {
  languages: [
    { name: "Go",         color: "#00acd7" },
    { name: "Java",       color: "#f59e0b" },
    { name: "TypeScript", color: "#3178c6" },
    { name: "SQL",        color: "#555"    },
  ],
  backend: [
    { name: "Kafka",      color: "#10b981" },
    { name: "PostgreSQL", color: "#336791" },
    { name: "Redis",      color: "#dc382d" },
    { name: "MongoDB",    color: "#47a248" },
  ],
  infra: [
    { name: "Kubernetes", color: "#326ce5" },
    { name: "Docker",     color: "#2496ed" },
    { name: "AWS",        color: "#ff9900" },
    { name: "Linux",      color: "#888"    },
  ],
  concepts: [
    { name: "Microservices",     color: "#8b5cf6" },
    { name: "System Design",     color: "#8b5cf6" },
    { name: "Distributed Sys",   color: "#8b5cf6" },
    { name: "Container Internals", color: "#8b5cf6" },
  ],
};

function fmt(d: string | null) {
  if (!d) return "";
  return new Date(d).toLocaleDateString("en-US", {
    month: "short", day: "numeric", year: "numeric",
  });
}

export default async function Home() {
  const [posts, projects] = await Promise.all([
    getPosts().catch(() => []),
    getProjects().catch(() => []),
  ]);

  const recent   = (posts ?? []).slice(0, 3);
  const featured = recent[0];
  const rest     = recent.slice(1);

  return (
    <div className="pb-16">

      {/* ── HERO ─────────────────────────────────────────── */}
      <div className="grid grid-cols-[1fr_400px] gap-16 items-center
                      py-24 border-b border-b1">
        <div>
          {/* badge */}
          <div className="inline-flex items-center gap-2 px-3.5 py-1.5 rounded-full
                          border border-[rgba(139,92,246,0.3)] bg-[rgba(139,92,246,0.08)]
                          text-[12px] text-p2 font-mono mb-7 tracking-[0.3px]">
            <span className="w-1.5 h-1.5 rounded-full bg-[var(--green)] animate-blink" />
            backend engineer · super.money
          </div>

          {/* title */}
          <h1 className="font-serif text-[60px] leading-[1.05] tracking-[-1.5px]
                         text-t1 mb-5">
            Ajinkya Singh<br />
            <span className="text-p2 italic">builds systems</span><br />
            that scale
          </h1>

          <p className="text-[17px] text-t2 leading-[1.8] max-w-[460px] mb-9 font-light">
            I work on distributed backend infrastructure at Super.money (Flipkart).
            I write deep technical posts about Go, container internals, Kafka, and
            Kubernetes — things I actually build at work.
          </p>

          <div className="flex gap-3 flex-wrap">
            <Link href="/blog"
                  className="inline-flex items-center gap-2 px-6 py-3 rounded-xl
                             bg-[var(--p)] text-white text-[15px] font-medium
                             hover:bg-[var(--p2)] transition-all hover:-translate-y-0.5">
              Read the blog
              <svg className="w-4 h-4" viewBox="0 0 16 16" fill="none">
                <path d="M3 8h10M9 4l4 4-4 4" stroke="currentColor"
                      strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round"/>
              </svg>
            </Link>
            <a href="https://github.com/a-ZINC" target="_blank"
               className="inline-flex items-center gap-2 px-6 py-3 rounded-xl
                          border border-b3 text-t2 text-[15px]
                          hover:text-t1 hover:bg-s1 transition-all">
              <svg className="w-4 h-4" viewBox="0 0 16 16" fill="currentColor">
                <path d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.013 8.013 0 0016 8c0-4.42-3.58-8-8-8z"/>
              </svg>
              github
            </a>
            <a href="https://dev.to/azinc" target="_blank"
               className="inline-flex items-center gap-2 px-6 py-3 rounded-xl
                          border border-b3 text-t2 text-[15px]
                          hover:text-t1 hover:bg-s1 transition-all">
              dev.to
            </a>
          </div>

          {/* stats */}
          <div className="flex gap-8 mt-12 pt-8 border-t border-b1">
            {[
              { n: "8.2k", l: "total reads"   },
              { n: "14",   l: "posts written"  },
              { n: "3",    l: "open source"    },
            ].map((s) => (
              <div key={s.l}>
                <div className="font-serif text-[26px] text-t1 tracking-[-0.5px]">
                  {s.n}
                </div>
                <div className="text-[12px] text-t3 mt-0.5 tracking-[0.3px]">
                  {s.l}
                </div>
              </div>
            ))}
          </div>
        </div>

        {/* code card */}
        <div className="bg-s1 border border-b2 rounded-2xl overflow-hidden">
          <div className="flex items-center gap-2 px-4 py-3.5 bg-s2 border-b border-b1">
            <span className="w-2.5 h-2.5 rounded-full bg-[#ff5f57]" />
            <span className="w-2.5 h-2.5 rounded-full bg-[#febc2e]" />
            <span className="w-2.5 h-2.5 rounded-full bg-[#28c840]" />
            <span className="ml-auto font-mono text-[11px] text-t3">
              conti/runtime.go
            </span>
          </div>
          <pre className="p-6 font-mono text-[13px] leading-[1.8] overflow-x-auto">
{`<span style="color:#524e7a">// container runtime in Go</span>

<span style="color:#c084fc">func </span><span style="color:#67e8f9">Run</span><span style="color:#f0eeff">(cfg Config) </span><span style="color:#c084fc">error </span><span style="color:#f0eeff">{</span>
  ns <span style="color:#f0eeff">:= </span><span style="color:#67e8f9">newNamespace</span><span style="color:#f0eeff">(</span>
    syscall.<span style="color:#fbbf24">CLONE_NEWPID</span><span style="color:#f0eeff">,</span>
    syscall.<span style="color:#fbbf24">CLONE_NEWNET</span><span style="color:#f0eeff">,</span>
    syscall.<span style="color:#fbbf24">CLONE_NEWNS</span><span style="color:#f0eeff">,</span>
  <span style="color:#f0eeff">)</span>

  <span style="color:#c084fc">if </span>err <span style="color:#f0eeff">:= </span>ns.<span style="color:#67e8f9">Apply</span><span style="color:#f0eeff">(cfg); </span>err <span style="color:#f0eeff">!= </span><span style="color:#c084fc">nil </span><span style="color:#f0eeff">{</span>
    <span style="color:#c084fc">return </span>err
  <span style="color:#f0eeff">}</span>

  <span style="color:#c084fc">return </span><span style="color:#67e8f9">setupCgroups</span><span style="color:#f0eeff">(cfg)</span>
<span style="color:#f0eeff">}</span>`}
          </pre>
        </div>
      </div>

      {/* ── POSTS ────────────────────────────────────────── */}
      <section className="mt-20">
        <div className="flex items-end justify-between mb-8">
          <div>
            <div className="text-[11px] tracking-[2.5px] uppercase text-t3
                            font-mono mb-2">writing</div>
            <h2 className="font-serif text-[34px] tracking-[-0.5px] text-t1">
              Recent posts
            </h2>
          </div>
          <Link href="/blog" className="text-[14px] text-p2 hover:text-p3
                                        transition-colors flex items-center gap-1">
            View all posts →
          </Link>
        </div>

        <div className="grid grid-cols-2 gap-4">

          {/* featured post */}
          {featured && (
            <Link href={`/blog/${featured.Slug}`}
                  className="col-span-2 grid grid-cols-2 gap-8 bg-s1 border border-b1
                             rounded-2xl p-7 hover:border-b2 hover:bg-s2
                             transition-all duration-300 hover:-translate-y-0.5 group">
              <div>
                {featured.SeriesName && (
                  <div className="inline-flex items-center gap-1.5 text-[11px]
                                  text-p2 bg-[rgba(139,92,246,0.1)]
                                  border border-[rgba(139,92,246,0.2)]
                                  px-3 py-1 rounded-full font-mono mb-3">
                    {featured.SeriesName}
                    {featured.SeriesPos && ` · part ${featured.SeriesPos}`}
                  </div>
                )}
                <div className="flex gap-1.5 mb-4 flex-wrap">
                  {featured.Tags?.map((t) => <Tag key={t} name={t} />)}
                </div>
                <h3 className="font-serif text-[26px] leading-[1.2] tracking-[-0.3px]
                               text-t1 mb-3 group-hover:text-p3 transition-colors">
                  {featured.Title}
                </h3>
                <p className="text-[14px] text-t2 leading-[1.7] mb-6">
                  {featured.Excerpt}
                </p>
                <div className="flex items-center gap-4 text-[12px] text-t3 font-mono">
                  <span>{fmt(featured.PublishedAt)}</span>
                  <span>{featured.ReadingTime} min read</span>
                  <span className="ml-auto text-p2 opacity-0
                                   group-hover:opacity-100 transition-opacity">
                    read →
                  </span>
                </div>
              </div>
              <div className="bg-s3 rounded-xl border border-b1
                              flex items-center justify-center">
                <svg className="w-12 h-12 opacity-20" viewBox="0 0 48 48" fill="none">
                  <rect x="4" y="4" width="18" height="18" rx="4"
                        stroke="currentColor" strokeWidth="2"/>
                  <rect x="26" y="4" width="18" height="18" rx="4"
                        stroke="currentColor" strokeWidth="2"/>
                  <rect x="4" y="26" width="18" height="18" rx="4"
                        stroke="currentColor" strokeWidth="2"/>
                  <rect x="26" y="26" width="18" height="18" rx="4"
                        stroke="currentColor" strokeWidth="2"/>
                </svg>
              </div>
            </Link>
          )}

          {/* rest of posts */}
          {rest.map((p) => (
            <Link key={p.ID} href={`/blog/${p.Slug}`}
                  className="bg-s1 border border-b1 rounded-2xl p-7
                             hover:border-b2 hover:bg-s2 transition-all
                             duration-300 hover:-translate-y-0.5 group">
              <div className="flex gap-1.5 mb-4 flex-wrap">
                {p.Tags?.map((t) => <Tag key={t} name={t} />)}
              </div>
              <h3 className="font-serif text-[21px] leading-[1.25] tracking-[-0.2px]
                             text-t1 mb-3 group-hover:text-p3 transition-colors">
                {p.Title}
              </h3>
              <p className="text-[13px] text-t2 leading-[1.7] mb-6 line-clamp-2">
                {p.Excerpt}
              </p>
              <div className="flex items-center gap-4 text-[12px] text-t3 font-mono">
                <span>{fmt(p.PublishedAt)}</span>
                <span>{p.ReadingTime} min read</span>
                <span className="ml-auto text-p2 opacity-0
                                 group-hover:opacity-100 transition-opacity">
                  →
                </span>
              </div>
            </Link>
          ))}

        </div>
      </section>

      {/* ── STACK ────────────────────────────────────────── */}
      <section className="mt-20 bg-s1 border border-b1 rounded-3xl p-12">
        <div className="text-[11px] tracking-[2.5px] uppercase text-t3
                        font-mono mb-2">tech stack</div>
        <h2 className="font-serif text-[28px] tracking-[-0.5px] text-t1 mb-10">
          What I work with
        </h2>
        <div className="grid grid-cols-4 gap-8">
          {Object.entries(stack).map(([group, items]) => (
            <div key={group}>
              <div className="text-[11px] tracking-[1.5px] uppercase text-t3
                              font-mono mb-4">{group}</div>
              <div className="flex flex-col gap-2">
                {items.map((item) => (
                  <div key={item.name}
                       className="flex items-center gap-2.5 px-3.5 py-2.5
                                  bg-s2 rounded-xl border border-b1 text-[13px]
                                  text-t2 hover:text-t1 hover:border-b2
                                  hover:bg-s3 transition-all cursor-default">
                    <span className="w-2 h-2 rounded-full flex-shrink-0"
                          style={{ background: item.color }} />
                    {item.name}
                  </div>
                ))}
              </div>
            </div>
          ))}
        </div>
      </section>

      {/* ── PROJECTS ─────────────────────────────────────── */}
      {projects && projects.length > 0 && (
        <section className="mt-20">
          <div className="flex items-end justify-between mb-8">
            <div>
              <div className="text-[11px] tracking-[2.5px] uppercase text-t3
                              font-mono mb-2">open source</div>
              <h2 className="font-serif text-[34px] tracking-[-0.5px] text-t1">
                Projects
              </h2>
            </div>
            <Link href="/projects"
                  className="text-[14px] text-p2 hover:text-p3
                             transition-colors">
              All projects →
            </Link>
          </div>
          <div className="grid grid-cols-3 gap-4">
            {projects.slice(0, 3).map((pr) => (
              <a key={pr.ID}
                 href={pr.RepoURL || pr.URL || "#"}
                 target="_blank" rel="noopener noreferrer"
                 className="bg-s1 border border-b1 rounded-2xl p-6
                            hover:border-b2 hover:bg-s2 transition-all
                            hover:-translate-y-0.5 group">
                <div className="flex items-center justify-between mb-4">
                  <div className="flex items-center gap-2 text-[12px]
                                  text-t3 font-mono">
                    <span className="w-2 h-2 rounded-full"
                          style={{
                            background:
                              pr.Language === "Go"         ? "#00acd7" :
                              pr.Language === "TypeScript" ? "#3178c6" :
                              pr.Language === "Java"       ? "#f59e0b" : "#888",
                          }} />
                    {pr.Language}
                  </div>
                  <span className="text-t3 group-hover:text-p2
                                   group-hover:translate-x-0.5 group-hover:-translate-y-0.5
                                   transition-all text-[16px]">↗</span>
                </div>
                <div className="text-[17px] font-medium text-t1 mb-2 tracking-[-0.2px]">
                  {pr.Name}
                </div>
                <div className="text-[13px] text-t2 leading-[1.65]">
                  {pr.Description}
                </div>
              </a>
            ))}
          </div>
        </section>
      )}

    </div>
  );
}