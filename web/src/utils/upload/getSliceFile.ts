import { fileHash } from "./fileManipulation";

export /**
* @param {File} file slice file
*@param {number} pieceSizes slice size
*@param {string} fileKey file unique identifier
*/
    const getSliceFile = async (file: File, pieceSizes = 1, fileKey: string) => {
        const piece = 1024 * 1024 * pieceSizes;
        //Total file size
        const totalSize = file.size;
        const fileName = file.name;
        //Start byte of each upload
        let start = 0;
        let index = 1;
        //End byte of each upload
        let end = start + piece;
        const chunks = [];
        while (start < totalSize) {
            const current = Math.min(end, totalSize);
            //Intercept the data that needs to be uploaded each time based on the length
            //File object inherits from Blob object, so it contains slice method
            const blob = file.slice(start, current);
const hash = (await fileHash(blob)) as string;

            //Special processing, docking with Alibaba Cloud to upload large files
            chunks.push({
                file: blob,
                size: totalSize,
                index,
                fileSizeInByte: totalSize,
                name: fileName,
                fileName,
                hash,
                sliceSizeInByte: blob.size,
                progress: 0,
                fileKey,
            });
            start = current;
            end = start + piece;
            index += 1;
        }
        return chunks;
    };