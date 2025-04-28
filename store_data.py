import os
from langchain_community.vectorstores import Chroma
from langchain_huggingface import HuggingFaceEmbeddings
from langchain.text_splitter import RecursiveCharacterTextSplitter
from load_files import load_repo_files

# Load documents from your repo
docs = load_repo_files("/Users/barkhasingh/go_test/golang-grpc-realworld-example")
print(f"\n Loaded {len(docs)} raw documents")

if not docs:
    raise ValueError(" No documents found! Check your load_files.py logic.")

# Split documents
text_splitter = RecursiveCharacterTextSplitter(
    chunk_size=1500,
    chunk_overlap=200,
    separators=["\n\n", "\n", " ", ""]
)
docs_split = text_splitter.split_documents(docs)
print(f" Split into {len(docs_split)} chunks")

# Extract content & metadata
texts = [doc.page_content for doc in docs_split]
metadatas = [doc.metadata for doc in docs_split]

# Embedding model
embedding_model = HuggingFaceEmbeddings(model_name="sentence-transformers/all-MiniLM-L6-v2")

# Store in Chroma
vectorstore = Chroma.from_texts(
    texts=texts,
    embedding=embedding_model,
    metadatas=metadatas,
    persist_directory="chroma_db"
)
vectorstore.persist()
print(" All documents embedded & stored in ChromaDB!")
