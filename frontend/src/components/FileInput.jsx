const FileInput = () => {
    
    return (
        <div className="rounded-md border border-gray-100 bg-white p-2 shadow-md w-1/12 mx-auto mt-3">
            <label htmlFor="upload" className="flex flex-row items-center gap-2 cursor-pointer text-center justify-center">
            <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6 fill-white stroke-btn" viewBox="0 0 24 24" stroke="currentColor" strokeWidth="2">
                <path strokeLinecap="round" strokeLinejoin="round" d="M9 13h6m-3-3v6m5 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
            <span className="text-btn font-semibold">Subir archivo</span>
            </label>
            <input id="upload" type="file" className="hidden" />
        </div>
    )
    
}
    
export default FileInput;